package models

func (ms *ModelSuite) Test_Group() {
	super := &Super{}
	group := Group{}
	group.Name = "That Group of Supers"
	super.Groups = append(super.Groups, group)
	_, err := ms.DB.ValidateAndCreate(super)
	if err != nil {
		ms.Fail("Erro ao validar e criar super com group")
	}
}
