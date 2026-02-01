package network

import (
	"net/http"
	"server-a/server/constant"
	"server-a/server/dto"

	"github.com/gin-gonic/gin"
)

func emailRouter(n *Network) {
	n.Router(POST, "/email/create", n.createMemberByEmail)
	n.Router(POST, "/email/login", n.loginWithEmail)
	n.Router(GET, "/email/check", n.checkEmail)
	n.Router(POST, "/email/otp/send", n.sendEmailOTP)
	n.Router(POST, "/email/otp/verify", n.verifyEmailOTP)
	n.Router(POST, "/email/apple", n.signInWithApple)
}

func (n *Network) createMemberByEmail(c *gin.Context) {
	var req dto.EmailMemberSaveRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	result, err := n.service.CreateMemberByEmail(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (n *Network) loginWithEmail(c *gin.Context) {
	var req dto.EmailMemberLoginRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	result, rt, err := n.service.LoginWithEmail(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	if rt != "" {
		c.SetCookie("refresh_token",
			rt,
			constant.RefreshTokenTTL,
			"",
			"",
			false,
			true,
		)
	}
	c.JSON(http.StatusOK, result)
}

func (n *Network) checkEmail(c *gin.Context) {
	var req dto.EmailRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx := c.Request.Context()
	ok, err := n.service.IsEmailUsable(ctx, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, ok)
}

func (n *Network) sendEmailOTP(c *gin.Context) {
	var req dto.EmailOTPSendRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	result, err := n.service.SendEmailOTP(c.Request.Context(), req.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (n *Network) verifyEmailOTP(c *gin.Context) {
	var req dto.EmailOTPVerifyRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	result, err := n.service.VerifyEmailOTP(req.OTP, req.VerificationId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (n *Network) signInWithApple(c *gin.Context) {
	var req dto.SignInWithAppleRequest
	ctx := c.Request.Context()
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	responseBody, err := n.service.SignInWithApple(
		ctx,
		req.User,
		req.Nonce,
		req.IdentityToken,
		req.Email,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
	}
	c.JSON(http.StatusOK, responseBody)
}
