package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	superOnDBFails := 0
	//Para cada caso de teste
	for i, tcase := range tcases {
		//Fazer chamada a api
		res := as.JSON(fmt.Sprintf("/%s/add?%s", os.Getenv("SUPER_API_KEY"), tcase.param)).Get()
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
		//Confere se os supers foram criados no banco de dados
		//Variavel para argumentos de where
		var supersToFindDB []interface{}
		supersToFindDB = append(supersToFindDB, tcase.containFullName)
		//Variavel para supers encontrados
		supersFoundOnDB := []models.Super{}
		//Executa query
		as.DB.Where("full_name in (?)", supersToFindDB...).All(&supersFoundOnDB)
		//Para cada resultado esperado
		for _, expectedFullName := range tcase.containFullName {
			superFound := false
			//Para cada resultado da query
			for _, super := range supersFoundOnDB {
				if super.FullName == expectedFullName {
					superFound = true
					break
				}
			}
			assertion := as.Assert().True(superFound, `fullname "%s" não encontrado no caso de teste %d`, expectedFullName, i)
			if !assertion {
				superOnDBFails++
			}
		}
	}
	//Se houver qualquer falha, falhar o teste
	as.Require().Equal(true, containFullNameFails < 1 && notContainFullNameFails < 1 && superOnDBFails < 1, fmt.Sprintf(`Falhas de containFullName:%d Falhas de notContainFullName:%d Falhas de superOnDBFails:%d`, containFullNameFails, notContainFullNameFails, superOnDBFails))
}

func (as *ActionSuite) Test_Supers_Search() {
	//Carrega fixtures
	as.LoadFixture("Bat Family")
	as.LoadFixture("Batman Villains")
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
		{"name=batman", http.StatusOK, "name", []string{"batman"}, []string{"batgirl", "batwoman"}},
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
				case "name":
					expectedName := strings.ToLower(expectedResult)
					superName := strings.ToLower(super.Name)
					assertion := as.Assert().Contains(superName, expectedName, fmt.Sprintf(`Todos os resultados do caso de teste %d deveriam conter "%s" no seu name. Name encontrado:"%s"`, i, expectedName, super.Name))
					if !assertion {
						allContainValueFails++
					}
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
				case "name":
					notExpectedName := strings.ToLower(notExpectedResult)
					superName := strings.ToLower(super.Name)
					assertion := as.Assert().NotContains(superName, notExpectedName, fmt.Sprintf(`Todos os resultados do caso de teste %d não deveriam conter "%s" no seu name. Name encontrado:"%s"`, i, notExpectedName, super.Name))
					if !assertion {
						allDontContainValueFails++
					}
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

func (as *ActionSuite) Test_Supers_Destroy() {
	//Carrega fixtures
	as.LoadFixture("Bat Family")
	//Carrega supers do banco de dados
	supers := []models.Super{}
	as.DB.All(&supers)
	//Para cada super, executar ação de destruição
	for _, super := range supers {
		//Chamada da rota
		res := as.JSON(fmt.Sprintf("/%s/destroy?super_id=%s", os.Getenv("SUPER_API_KEY"), super.ID)).Get()
		as.Require().Equal(http.StatusOK, res.Code, fmt.Sprintf("Codigo de Resposta esperado: %d Codigo de Resposta recebido: %d", http.StatusOK, res.Code))
	}
	//Confere se o banco de dados esta vazio
	supers = []models.Super{}
	as.DB.All(&supers)
	as.Require().Equal(0, len(supers))
}

func (as *ActionSuite) Test_Supers_Heros() {
	//Carrega fixtures
	as.LoadFixture("Bat Family")
	as.LoadFixture("Batman Villains")
	//Constante com numero de herois nas fixtures
	const heroNumber = 10
	//Executa chamada a rota
	res := as.JSON("/heros").Get()
	//Confere codigo de resposta
	as.Require().Equal(http.StatusOK, res.Code, fmt.Sprintf("Codigo de resposta inesperado: %d Esperava: %d", res.Code, http.StatusOK))
	//Carrega supers do json
	supers := []models.Super{}
	json.Unmarshal(res.Body.Bytes(), &supers)
	//Confere numero de herois
	as.Assert().Equal(heroNumber, len(supers), fmt.Sprintf("Esperava encontrar %d supers, mas encontrou %d", heroNumber, len(supers)))
	//Para cada super, testa de aligment é good
	for _, super := range supers {
		as.Require().Equal("good", strings.ToLower(super.Alignment), fmt.Sprintf(`Essa ação só deveria retornar supers com o aligment "good", mas retornou o super "%s" que possui aligment "%s"`, super.Name, strings.ToLower(super.Alignment)))
	}
}

func (as *ActionSuite) Test_Supers_Villains() {
	//Carrega fixtures
	as.LoadFixture("Bat Family")
	as.LoadFixture("Batman Villains")
	//Constante com numero de vilões nas fixtures
	const villainNumber = 3
	//Executa chamada a rota
	res := as.JSON("/villains").Get()
	//Confere codigo de resposta
	as.Require().Equal(http.StatusOK, res.Code, fmt.Sprintf("Codigo de resposta inesperado: %d Esperava: %d", res.Code, http.StatusOK))
	//Carrega supers do json
	supers := []models.Super{}
	json.Unmarshal(res.Body.Bytes(), &supers)
	//Confere numero de resultados
	as.Assert().Equal(villainNumber, len(supers), fmt.Sprintf("Esperava encontrar %d supers, mas encontrou %d", villainNumber, len(supers)))
	//Para cada super, testa de aligment é bad
	for _, super := range supers {
		as.Require().Equal("bad", strings.ToLower(super.Alignment), fmt.Sprintf(`Essa ação só deveria retornar supers com o aligment "bad", mas retornou o super "%s" que possui aligment "%s"`, super.Name, strings.ToLower(super.Alignment)))
	}
}

func (as *ActionSuite) Test_Supers_All() {
	//Carrega fixtures
	as.LoadFixture("Bat Family")
	as.LoadFixture("Batman Villains")
	//Constante com numero de vilões nas fixtures
	const heroNumber = 10
	const villainNumber = 3
	const totalNumber = heroNumber + villainNumber
	//Executa chamada a rota
	res := as.JSON("/all").Get()
	//Confere codigo de resposta
	as.Require().Equal(http.StatusOK, res.Code, fmt.Sprintf("Codigo de resposta inesperado: %d Esperava: %d", res.Code, http.StatusOK))
	//Carrega supers do json
	supers := []models.Super{}
	json.Unmarshal(res.Body.Bytes(), &supers)
	//Confere numero de resultados
	as.Assert().Equal(totalNumber, len(supers), fmt.Sprintf("Esperava encontrar %d supers, mas encontrou %d", totalNumber, len(supers)))
}
