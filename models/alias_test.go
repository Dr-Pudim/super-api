package models

func (ms *ModelSuite) Test_Alias() {
	super := &Super{}
	alias := Alias{}
	alias.Name = "The Other Name"
	super.Aliases = append(super.Aliases, alias)
	_, err := ms.DB.ValidateAndCreate(super)
	if err != nil {
		ms.Fail("Erro ao validar e criar super com alias")
	}
}
