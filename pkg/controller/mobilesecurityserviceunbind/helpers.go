package mobilesecurityserviceunbind

import (
	"github.com/aerogear/mobile-security-service-operator/pkg/models"
)

//hasApp return true when APP has ID which is just created by the REST Service API
func hasApp(app models.App) bool {
	return len(app.ID) > 0
}