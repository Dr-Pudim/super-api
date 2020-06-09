package models

func (ms *ModelSuite) Test_Relative() {
	super := &Super{}
	relative := Relative{}
	relative.Name = "Relative (Relation)"
	super.Relatives = append(super.Relatives, relative)
	_, err := ms.DB.ValidateAndCreate(super)
	if err != nil {
		ms.Fail("Erro ao validar e criar super")
	}
}
