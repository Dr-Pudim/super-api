package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
			assertion := as.Assert().Equal(true, fullNameFound, fmt.Sprintf(`Caso de teste %d deveria conter super com "fullname":"%s"`, i, name))
			//Se o nome não foi encontrado, acrescentar uma falha
			if !assertion {
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
			assertion := as.Assert().Equal(false, fullNameFound, fmt.Sprintf(`Caso de teste %d não deveria conter super com "fullname":"%s"`, i, name))
			if !assertion {
				notContainFullNameFails++
			}
		}
	}
	//Se houver qualquer falha, falhar o teste
	as.Require().Equal(true, containFullNameFails < 1 && notContainFullNameFails < 1, fmt.Sprintf(`Falhas de containFullName:%d Falhas de notContainFullName:%d`, containFullNameFails, notContainFullNameFails))
}

func (as *ActionSuite) Test_Supers_Search() {
	//Carrega fixtures
	as.LoadFixture("Bat Family")
	//Struct de casos de teste
	tcases := []struct {
		searchParams        string
		expectedRespCode    int
		fieldToTest         string
		allContainValue     []string
		allDontContainValue []string
	}{
		{"gender=male", http.StatusOK, "gender", []string{"male"}, []string{"female"}},
		{"gender=female", http.StatusOK, "gender", []string{"female"}, []string{"male"}},
		{"intelligence=70", http.StatusOK, "intelligence", []string{"70"}, []string{}},
	}
	//Variaveis para contar falhas
	allContainValueFails := 0
	allDontContainValueFails := 0
	//Para cada caso de teste
	for i, tcase := range tcases {
		//Fazer chamada a api
		res := as.JSON("/search?" + tcase.searchParams).Get()
		//Conferir codigo de resposta
		as.Equal(tcase.expectedRespCode, res.Code, fmt.Sprintf(`Codigo de Resposta esperado:%d Codigo de Resposta recebido:%d`, tcase.expectedRespCode, res.Code))
		//Unmarshal do json
		response := []models.Super{}
		json.Unmarshal(res.Body.Bytes(), &response)
		//Confere se a resposta contem os resultados esperados
		//Para cada resultado esperado
		for _, expectedResult := range tcase.allContainValue {
			//Para cada super da resposta
			for _, super := range response {
				//Seleciona campo para testar
				switch tcase.fieldToTest {
				case "gender":
					lowerCaseGender := strings.ToLower(super.Gender)
					assertion := as.Assert().Equal(expectedResult, lowerCaseGender, fmt.Sprintf(`Todos os resutlado do caso de teste %d deveriam conter o campo "gender" com valor "%s"`, i, expectedResult))
					if !assertion {
						allContainValueFails++
					}
				case "intelligence":
					expectedIntelligence, err := strconv.Atoi(expectedResult)
					if err != nil {
						as.Fail(fmt.Sprintf(`Erro na conversão do valor do campo Intelligence do super: %s`, super.Name))
					}
					assertion := as.Assert().GreaterOrEqual(super.Intelligence, expectedIntelligence, fmt.Sprintf(`Todos os resultados do caso de teste %d deveriam conter intelligence iqual ou maior que %d`, i, expectedIntelligence))
					if !assertion {
						allContainValueFails++
					}
				}
			}
		}
		//Confere se a resposta contem os resultados *não* esperados
		//Para cada resultado *não* esperado
		for _, notExpectedResult := range tcase.allDontContainValue {
			//Para cada super da resposta
			for _, super := range response {
				//Seleciona campo para testar
				switch tcase.fieldToTest {
				case "gender":
					lowerCaseGender := strings.ToLower(super.Gender)
					assertion := as.Assert().NotEqual(notExpectedResult, lowerCaseGender, fmt.Sprintf(`Todos os resutlado do caso de teste %d *não* deveriam conter o campo "gender" com valor "%s"`, i, notExpectedResult))
					if !assertion {
						allDontContainValueFails++
					}
				}
			}
		}
	}
	//Se houver qualquer falha, falhar o teste
	as.Require().Equal(true, allContainValueFails < 1 && allDontContainValueFails < 1, fmt.Sprintf(`Falhas de allContainValueFails:%d Falhas de allDontContainValueFails:%d`, allContainValueFails, allDontContainValueFails))
}
