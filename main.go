package main

import (
	"fmt"
	"webauthn/config"
	"webauthn/database"
	"webauthn/model"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

func main() {
	config.InitConfig()
	database.InitDB()
	r := gin.Default()
	wconfig := &webauthn.Config{
		RPDisplayName: "Go Webauthn",                          // Display Name for your site
		RPID:          "webauthn.sylu.site",                   // Generally the FQDN for your site
		RPOrigins:     []string{"https://webauthn.sylu.site"}, // The origin URLs allowed for WebAuthn requests
	}

	_, err := webauthn.New(wconfig)
	if err != nil {
		fmt.Println(err)
	}
	r.GET("/get_register", GetRegister)
	r.POST("/register", Register)
	r.POST("/login", Login)
	r.GET("/get_login", GetLogin)
	r.Run(":8080")
}
func GetRegister(ctx *gin.Context) {
	name := ctx.Param("name")
	// challengeid, err := protocol.CreateChallenge()
	// if err != nil {
	// 	fmt.Println("get register err", err)
	// }
	DB := database.GetDB()
	var user model.User
	DB.Table("users").Where("name = ?", name).First(&user)
	var userid int64
	DB.Table("users").Count(&userid)
	userid += 1
	if user.ID == 0 {
		ctx.JSON(200, gin.H{
			"code": 200,
			"date": gin.H{
				// "challengeid": challengeid.String(),
				"challengeid": "1234567890abcdef",
				"userid":      userid,
			},
			"msg": "get register successed",
		})
	} else {
		ctx.JSON(200, gin.H{
			"code": 422,
			"date": "",
			"msg":  "username is in using",
		})
	}
}
func Register(ctx *gin.Context) {
	name := "name"
	var webauthnJson protocol.CredentialCreationResponse
	ctx.ShouldBindJSON(&webauthnJson)
	fmt.Println(webauthnJson)
	// var pcc *protocol.ParsedCredentialCreationData
	pcc, _ := webauthnJson.Parse()
	err := pcc.Verify("1234567890abcdef", false, "webauthn.sylu.site", []string{pcc.Response.CollectedClientData.Origin})
	if err != nil {
		fmt.Println("err", err)
	}
	credential, _ := webauthn.MakeNewCredential(pcc)
	DB := database.GetDB()
	newUser := model.User{
		Name:         name,
		DisplayName:  "displayName",
		Certificates: credential.ID,
	}
	DB.Table("users").Create(&newUser)
	DB.Table("users").Where("name = ?", name).First(&newUser)
	newCerd := model.Certificate{
		UserID:     newUser.ID,
		Credention: credential.PublicKey,
	}
	DB.Table("certificates").Create(&newCerd)
	ctx.JSON(200, gin.H{
		"code": "200",
		"data": "success",
		"msg":  "webauthn",
	})
}

func GetLogin(ctx *gin.Context) {
	DB := database.GetDB()
	name := "name"
	var user model.User
	DB.Table("users").Where("name = ?", name).First(&user)
	var json_data protocol.PublicKeyCredentialRequestOptions
	json_data.AllowedCredentials = []protocol.CredentialDescriptor{
		protocol.CredentialDescriptor{
			Type:            "public-key",
			CredentialID:    user.Certificates,
			Transport:       []protocol.AuthenticatorTransport{"usb", "nfc", "ble"},
			AttestationType: "",
		},
	}
	ctx.JSON(200, json_data)
}

func Login(ctx *gin.Context) {
	DB := database.GetDB()
	name := "name"
	var user model.User
	DB.Table("users").Where("name = ?", name).First(&user)
	var webauthnJson protocol.CredentialAssertionResponse
	ctx.ShouldBindJSON(&webauthnJson)
	fmt.Println(webauthnJson)
	pca, _ := webauthnJson.Parse()
	var cred model.Certificate

	DB.Table("certificates").Where("user_id = ?", user.ID).Scan(&cred)
	fmt.Println(cred.Credention)
	fmt.Println(pca.Verify("1234567890abcdef", "webauthn.sylu.site", []string{pca.Response.CollectedClientData.Origin}, "", false, cred.Credention))
	ctx.JSON(200, gin.H{
		"code": "200",
		"data": "success",
		"msg":  "webauthn",
	})
}
