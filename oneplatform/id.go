package oneplatform

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/patcharp/go_swth/requests"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

const IdProductionEndpoint = "https://one.th"

type Identity struct {
	ApiEndpoint  string
	ClientId     string
	ClientSecret string
}

type AuthenticationResult struct {
	TokenType    string         `json:"token_type"`
	ExpiresIn    int            `json:"expires_in"`
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	AccountID    string         `json:"account_id"`
	Result       string         `json:"result"`
	Username     string         `json:"username"`
	Profile      AccountProfile `json:"profile"`
}

type AccountProfile struct {
	ID                 string          `json:"id"`
	FirstNameTH        string          `json:"first_name_th"`
	LastNameTH         string          `json:"last_name_th"`
	FirstNameENG       string          `json:"first_name_eng"`
	LastNameENG        string          `json:"last_name_eng"`
	TitleTH            string          `json:"account_title_th"`
	TitleENG           string          `json:"account_title_eng"`
	IDCardType         string          `json:"id_card_type"`
	IDCardTypeNumber   string          `json:"id_card_num"`
	IDCardHashed       string          `json:"hash_id_card_num"`
	AccountCategory    string          `json:"account_category"`
	AccountSubCategory string          `json:"account_sub_category"`
	ThaiEmail1         string          `json:"thai_email"`
	ThaiEmail2         string          `json:"thai_email2"`
	StatusCD           string          `json:"status_cd"`
	BirthDate          string          `json:"birth_date"`
	StatusDate         string          `json:"status_dt"`
	RegisterDate       string          `json:"register_dt"`
	AddressID          string          `json:"address_id"`
	CreatedAt          string          `json:"created_at"`
	CreatedBy          string          `json:"created_by"`
	UpdatedAt          string          `json:"updated_at"`
	UpdatedBy          string          `json:"updated_by"`
	Reason             string          `json:"reason"`
	TelephoneNumber    string          `json:"tel_no"`
	NameOnDocTH        string          `json:"name_on_document_th"`
	NameOnDocENG       string          `json:"name_on_document_eng"`
	Mobile             []AccountMobile `json:"mobile"`
	Email              []AccountEmail  `json:"email"`
	Address            []string        `json:"address"`
	AccountAttr        []string        `json:"account_attribute"`
	Status             string          `json:"status"`
	LastUpdate         string          `json:"last_update"`
	Employee           *Employee       `json:"has_employee_detail"`
}

type AccountMobile struct {
	ID           string             `json:"id"`
	MobileNumber string             `json:"mobile_no"`
	CreatedAt    string             `json:"created_at"`
	CreatedBy    string             `json:"created_by"`
	UpdatedAt    string             `json:"updated_at"`
	UpdatedBy    string             `json:"updated_by"`
	DeletedAt    string             `json:"deleted_at"`
	MobilePivot  AccountMobilePivot `json:"pivot"`
}

type AccountMobilePivot struct {
	AccountID   string `json:"account_id"`
	MobileID    string `json:"mobile_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	StatusCD    string `json:"status_cd"`
	PrimaryFlag string `json:"primary_flg"`
	ConfirmFlag string `json:"mobile_confirm_flg"`
	ConfirmDate string `json:"mobile_confirm_dt"`
}

type AccountEmail struct {
	ID         string            `json:"id"`
	Email      string            `json:"email"`
	CreatedAt  string            `json:"created_at"`
	CreatedBy  string            `json:"created_by"`
	UpdatedAt  string            `json:"updated_at"`
	UpdatedBy  string            `json:"updated_by"`
	DeletedBy  string            `json:"deleted_at"`
	EmailPivot AccountEmailPivot `json:"pivot"`
}

type AccountEmailPivot struct {
	AccountID   string `json:"account_id"`
	EmailID     string `json:"email_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	StatusCD    string `json:"status_cd"`
	PrimaryFlag string `json:"primary_flg"`
	ConfirmFlag string `json:"email_confirm_flg"`
	ConfirmDate string `json:"email_confirm_dt"`
}

type Employee struct {
	Id         uuid.UUID       `json:"id"`
	AccountId  string          `json:"account_id"`
	BizId      string          `json:"biz_id"`
	Email      string          `json:"email"`
	EmployeeId string          `json:"employee_id"`
	Account    *AccountProfile `json:"account"`
	Employee   *Employee       `json:"employee"`
	Position   string          `json:"position"`
	PositionId uuid.UUID       `json:"role_id"`
}

func NewIdentity(clientId string, clientSecret string) Identity {
	return Identity{
		ApiEndpoint:  IdProductionEndpoint,
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
}

func (id *Identity) SetEndpoint(ep string) {
	id.ApiEndpoint = ep
}

func (id *Identity) Login(username string, password string, profile bool) (AuthenticationResult, error) {
	var result AuthenticationResult
	body, _ := json.Marshal(&struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Username     string `json:"username"`
		Password     string `json:"password"`
	}{
		ClientID:     id.ClientId,
		ClientSecret: id.ClientSecret,
		GrantType:    "password",
		Username:     username,
		Password:     password,
	})
	headers := map[string]string{
		echo.HeaderContentType: "application/json",
	}
	r, err := requests.Post(id.url("/api/oauth/getpwd"), headers, bytes.NewBuffer(body), 0)
	if err != nil {
		return result, err
	}
	if r.Code != http.StatusOK {
		return result, errors.New(fmt.Sprintf("client return error with code %d (%s)", r.Code, string(r.Body)))
	}
	if err := json.Unmarshal(r.Body, &result); err != nil {
		return result, err
	}
	if profile {
		result.Profile, err = id.profile(result.TokenType, result.AccessToken)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

func (id *Identity) RefreshNewToken(refreshToken string) (AuthenticationResult, error) {
	var result AuthenticationResult
	if refreshToken == "" {
		return result, errors.New("unauthorized identity")
	}
	body, _ := json.Marshal(&struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RefreshToken string `json:"refresh_token"`
	}{
		ClientID:     id.ClientId,
		ClientSecret: id.ClientSecret,
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	})
	headers := map[string]string{
		echo.HeaderContentType: "application/json",
	}
	resp, err := requests.Post(id.url("/api/oauth/get_refresh_token"), headers, bytes.NewBuffer(body), 0)
	if err != nil {
		return result, err
	}
	if resp.Code != http.StatusOK {
		return result, errors.New(fmt.Sprintf("client return error with code %d (%s)", resp.Code, string(resp.Body)))
	}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return result, err
	}
	return result, nil
}

func (id *Identity) profile(tokenType string, accessToken string) (AccountProfile, error) {
	var profile AccountProfile
	if tokenType == "" || accessToken == "" {
		return profile, errors.New("login required")
	}
	headers := map[string]string{
		echo.HeaderAuthorization: fmt.Sprintf("%s %s", tokenType, accessToken),
	}
	rawResponse, err := requests.Get(id.url("/api/account"), headers, nil, 0)
	if err != nil {
		return profile, err
	}
	return profile, json.Unmarshal(rawResponse.Body, &profile)
}

func (id *Identity) url(path string) string {
	return fmt.Sprintf("%s%s", id.ApiEndpoint, path)
}
