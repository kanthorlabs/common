package testdata

import "github.com/jaswdr/faker"

var Fake faker.Faker

func init() {
	Fake = faker.New()
}
