package Login

import "github.com/gofiber/fiber/v2"

func TwitterIcon(c *fiber.Ctx) error {

	c.Render("TwitterLogin", nil)
	return nil
}
