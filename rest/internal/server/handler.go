package server

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"strings"
)

type UserDto struct {
	FullName string `json:"fullname"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     int32  `json:"role"`
}
type CredentialDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s RestServer) Signup(c echo.Context) error {
	u := new(UserDto)
	if err := c.Bind(u); err != nil {
		return err
	}
	res, err := s.authGrpc.Signup(u.FullName, u.Email, u.UserName, u.Password, u.Role)
	if err != nil {
		fmt.Println(err)
		return c.JSON(401, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (s RestServer) Login(c echo.Context) error {
	dto := new(CredentialDto)
	if err := c.Bind(dto); err != nil {
		return err
	}
	res, err := s.authGrpc.Login(dto.Username, dto.Password)
	if err != nil {
		return echo.ErrUnauthorized
	}
	return c.JSON(http.StatusOK, res)
}

func (s RestServer) ValidateToken(c echo.Context) error {
	user := c.Get("user")
	if user == nil {
		return c.String(http.StatusUnauthorized, "wrong token")
	}
	return c.JSON(200, user)
}

func (s RestServer) CreateCategory(c echo.Context) error {
	dto := new(struct {
		Name        string
		Parent      string
		Description string
	})
	if err := c.Bind(dto); err != nil {
		c.Error(err)
	}
	res, err := s.productGrpc.CreateCategory(dto.Name, dto.Parent, dto.Description)
	if err != nil {
		c.Error(err)
	}
	return c.JSON(200, res)
}

func (s RestServer) GetCategories(c echo.Context) error {
	res, err := s.productGrpc.GetCategories()
	if err != nil {
		return c.String(400, err.Error())
	}
	return c.JSON(200, res)
}

func (s RestServer) CreateProduct(c echo.Context) error {
	body := c.Request().Body
	dto := new(struct {
		Name     string            `json:"name"`
		Category string            `json:"category"`
		ImageUrl string            `json:"image_url"`
		MinPrice int32             `json:"min_price"`
		Fields   map[string]string `json:"fields"`
	})
	err := json.NewDecoder(body).Decode(&dto)
	fmt.Println(dto.Fields)
	if err != nil {
		return c.String(http.StatusBadRequest, "wrong data format")
	}
	res, err := s.productGrpc.CreateProduct(dto.Name, dto.Category, dto.ImageUrl, dto.MinPrice, dto.Fields)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(200, res)
}

func (s RestServer) GetProduct(c echo.Context) error {
	url := c.Request().URL
	paths := strings.Split(url.String(), "/")
	id := paths[len(paths)-1]
	res, err := s.productGrpc.GetProduct(id)
	if err != nil {
		return c.String(404, err.Error())
	}
	fmt.Println(res)
	fields := make(map[string]string)
	err = mapstructure.Decode(res.Fields, &fields)
	fields["id"] = res.GetId()
	fields["name"] = res.GetName()
	fields["category"] = res.GetCategory()
	fields["image_url"] = res.GetImageUrl()
	fields["min_price"] = fmt.Sprintf("%d", res.GetMinPrice())
	if err != nil {
		return c.String(404, err.Error())
	}
	return c.JSONPretty(200, fields, " ")
}

func (s RestServer) GetProductsByType(c echo.Context) error {
	//dto := new(struct {
	//	Category string `json:"category"`
	//})
	//_ = c.Bind(dto)
	category := c.QueryParam("category")
	prods := make([]map[string]string, 0)
	products, err := s.productGrpc.GetProductsByType(category)
	for _, product := range products {
		prod := make(map[string]string)
		err = mapstructure.Decode(product.Fields, &prod)
		prod["id"] = product.GetId()
		prod["name"] = product.GetName()
		prod["category"] = product.GetCategory()
		prod["image_url"] = product.GetImageUrl()
		prod["min_price"] = fmt.Sprintf("%d", product.GetMinPrice())
		prods = append(prods, prod)
	}
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, prods)
}

func (s RestServer) AddStore(c echo.Context) error {
	dto := new(struct {
		OwnerId   string `json:"owner_id"`
		StoreName string `json:"store_name"`
		StoreUrl  string `json:"store_url"`
		City      string `json:"city"`
	})
	if err := c.Bind(dto); err != nil {
		return c.String(400, err.Error())
	}
	res, err := s.supplierGrpc.AddStore(dto.OwnerId, dto.StoreName, dto.StoreUrl, dto.City)
	if err != nil {
		return c.JSON(400, err.Error())
	}
	return c.String(http.StatusOK, res)
}

func (s RestServer) SubmitOffer(c echo.Context) error {
	dto := new(struct {
		StoreId   string `json:"store_id"`
		ProductId string `json:"product_id"`
		Url       string `json:"url"`
		Desc      string `json:"description"`
		Price     int    `json:"price"`
	})
	if err := c.Bind(dto); err != nil {
		return c.String(400, err.Error())
	}
	res, err := s.supplierGrpc.SubmitOffer(dto.StoreId, dto.ProductId, dto.Url, dto.Desc, int32(dto.Price))
	if err != nil {
		return c.String(400, err.Error())
	}
	return c.String(http.StatusOK, res)
}

func (s RestServer) GetProductOffers(c echo.Context) error {
	url := c.Request().URL
	paths := strings.Split(url.String(), "/")
	prodId := paths[len(paths)-1]
	offers, err := s.supplierGrpc.GetProductOffers(prodId)
	//var response []*struct {
	//	StoreName string `json:"store_name"`
	//	StoreCity string `json:"store_city"`
	//	Price int32 `json:"price"`
	//	ProdDesc string `json:"description"`
	//	Url string `json:"url"`
	//}
	fmt.Println("offers : ", offers)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, offers)
}

func (s RestServer) GetOwnerStores(c echo.Context) error {
	ownerId := c.QueryParam("ownerId")
	res, err := s.supplierGrpc.GetStores(ownerId)
	if err != nil {
		return c.String(400, err.Error())
	}
	return c.JSON(200, res)
}

func (s RestServer) GetStoreInfo(c echo.Context) error {
	url := c.Request().URL
	paths := strings.Split(url.String(), "/")
	storeId := paths[len(paths)-1]
	res, err := s.supplierGrpc.GetStoreInfo(storeId)
	if err != nil {
		return echo.NewHTTPError(400)
	}
	return c.JSON(200, res)
}
