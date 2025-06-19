package main

import "aigendrug.com/router-core/internal/interface/rest"

// @title           ATP Central Router Core API
// @version         1.0.0
// @description     ATP Central Router Core API
// @termsOfService  http://swagger.io/terms/

// @Schemes   https
// @host      api.apt.aigendrug.com
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	rest.Run()
}
