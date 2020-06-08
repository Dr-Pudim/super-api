package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"super_api/models"
)

func (as *ActionSuite) Test_Supers_Create() {
	//Struct de casos de teste
	tcases := []struct {
		param              string
		expectedRespCode   int
		containFullName    []string
		notContainFullName []string
	}{
		{"name=batman", http.StatusCreated, []string{"Terry McGinnis", "Bruce Wayne", "Dick Grayson"}, []string{"Barbara Gordon"}},
		{"name=spider-man", http.StatusCreated, []string{"Peter Parker", "Miguel O'Hara", "Miles Morales"}, []string{"Benjamin Reilly", "May 'Mayday' Parker", "Jessica Drew"}},
		{"", http.StatusBadRequest, []string{}, []string{}},
	}
	//Variaveis para contar falhas
	containFullNameFails := 0
	notContainFullNameFails := 0
	//Para cada caso de teste
	for i, tcase := range tcases {
		//Fazer chamada a api
		res := as.JSON("/add?" + tcase.param).Get()
		//Conferir codigo de resposta
		as.Equal(tcase.expectedRespCode, res.Code, fmt.Sprintf(`Codigo de Resposta esperado:%d Codigo de Resposta recebido:%d`, tcase.expectedRespCode, res.Code))
		//Unmarshal do json
		response := []models.Super{}
		json.Unmarshal(res.Body.Bytes(), &response)
		//Confere se a resposta contem os nomes esperados
		for _, name := range tcase.containFullName {
			fullNameFound := false
			for _, super := range response {
				if super.FullName == name {
					fullNameFound = true
					break
				}
			}
			//Testa se o nome foi encontrado
			as.Assert().Equal(true, fullNameFound, fmt.Sprintf(`Caso de teste %d deveria conter super com "fullname":"%s"`, i, name))
			//Se o nome não foi encontrado, acrescentar uma falha
			if !fullNameFound {
				containFullNameFails++
			}
		}
		//Confere se a resposta *não* contem os nomes que se esperam que não tenha
		for _, name := range tcase.notContainFullName {
			fullNameFound := false
			for _, super := range response {
				if super.FullName == name {
					fullNameFound = true
				}
			}
			as.Assert().Equal(false, fullNameFound, fmt.Sprintf(`Caso de teste %d não deveria conter super com "fullname":"%s"`, i, name))
			if fullNameFound {
				notContainFullNameFails++
			}
		}
	}
	//Se houver qualquer falha, falhar o teste
	as.Require().Equal(true, containFullNameFails < 1 && notContainFullNameFails < 1, fmt.Sprintf(`Falhas de containFullName:%d Falhas de notContainFullName:%d`, containFullNameFails, notContainFullNameFails))
}
