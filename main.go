package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
)

func main() {
	r := gin.Default()
	r.POST("/", GetPostParmars)
	r.Run(":3000")
}

type WebAuthnRegister struct {
	Username string `json:"username"`
	Response struct {
		ID       string `json:"id"`
		RawID    string `json:"rawId"`
		Response struct {
			AuthenticatorData string `json:"authenticatorData"`
			ClientDataJSON    string `json:"clientDataJSON"`
			Signature         string `json:"signature"`
		} `json:"response"`
		Type                   string `json:"type"`
		ClientExtensionResults struct {
		} `json:"clientExtensionResults"`
		AuthenticatorAttachment interface{} `json:"authenticatorAttachment"`
	} `json:"response"`
}

func GetPostParmars(ctx *gin.Context) {
	var webauthnJson WebAuthnRegister
	protocol.ParseCredentialCreationResponseBody()
	ctx.ShouldBindJSON(&webauthnJson)
	ctx.JSON(200, gin.H{
		"code": "200",
		"data": "success",
		"msg":  "webauthn",
	})
}
