package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"super_api/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
)

func isNullValue(s string) bool {
	if s == "null" || s == "-" || s == "" {
		return true
	}
	return false
}

func convertFirstWordToInt(s string) (int, error) {
	i := strings.Index(s, " ")
	if i > -1 {
		return strconv.Atoi(s[:i])
	}
	return strconv.Atoi(s)
}

func stringContainsSubstring(s string, subString string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(subString))
}

func filterSupersName(supers []models.Super, name string) []models.Super {
	//Aloca array de supers
	filteredSupers := []models.Super{}
	//Para cada super, comparar name
	for _, super := range supers {
		if stringContainsSubstring(super.Name, name) {
			filteredSupers = append(filteredSupers, super)
		}
	}
	return filteredSupers
}

func filterSupersFullName(supers []models.Super, name string) []models.Super {
	//Aloca array de supers
	filteredSupers := []models.Super{}
	//Para cada super, comparar name
	for _, super := range supers {
		if stringContainsSubstring(super.FullName, name) {
			filteredSupers = append(filteredSupers, super)
		}
	}
	return filteredSupers
}

func filterSupersOccupation(supers []models.Super, name string) []models.Super {
	//Aloca array de supers
	filteredSupers := []models.Super{}
	//Para cada super, comparar name
	for _, super := range supers {
		if stringContainsSubstring(super.Occupation, name) {
			filteredSupers = append(filteredSupers, super)
		}
	}
	return filteredSupers
}

type powerStats struct {
	Intelligence string `json:"intelligence"`
	Strength     string `json:"strength"`
	Speed        string `json:"speed"`
	Durability   string `json:"durability"`
	Power        string `json:"power"`
	Combat       string `json:"combat"`
}

type biography struct {
	FullName        string   `json:"full-name"`
	AlterEgos       string   `json:"alter-egos"`
	Aliases         []string `json:"aliases"`
	PlaceOfBirth    string   `json:"place-of-birth"`
	FirstAppearance string   `json:"first-appearance"`
	Publisher       string   `json:"publisher"`
	Alignment       string   `json:"alignment"`
}

type appearance struct {
	Gender    string   `json:"gender"`
	Race      string   `json:"race"`
	Height    []string `json:"height"`
	Weight    []string `json:"weight"`
	EyeColor  string   `json:"eye-color"`
	HairColor string   `json:"hair-color"`
}

type work struct {
	Occupation string `json:"occupation"`
	Base       string `json:"base"`
}

type connections struct {
	GroupAffiliation string `json:"group-affiliation"`
	Relatives        string `json:"relatives"`
}

type image struct {
	URL string `json:"url"`
}

type character struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Powerstats  powerStats  `json:"powerstats"`
	Biography   biography   `json:"biography"`
	Appearance  appearance  `json:"appearance"`
	Work        work        `json:"work"`
	Connections connections `json:"connections"`
	Image       image       `json:"image"`
}

//SearchResponse é uma struct que representa uma resposta de pesquisa da superheroapi.com
type SearchResponse struct {
	Response   string      `json:"response"`
	Resultsfor string      `json:"results-for"`
	Results    []character `json:"results"`
}

// SupersCreate default implementation.
func SupersCreate(c buffalo.Context) error {
	//Lê nome no parametro da rota
	param := c.Param("name")
	//Se não houver o parametro name, retornar mensagem
	if param == "" {
		return c.Render(http.StatusOK, r.JSON(map[string]string{"message": "Sem parametro name para buscar"}))
	}
	//Gera url de pesquisa para consultar a superheroapi
	url := fmt.Sprintf("https://superheroapi.com/api/%s/search/%s", os.Getenv("SUPERHEROAPI_ACCESS_TOKEN"), param)
	//Faz chama a superheroapi para procurar o super
	resp, err := http.Get(url)
	//Tratamento de erro da chamada
	if err != nil {
		return c.Render(http.StatusOK, r.JSON(map[string]string{"message": "Erro na chamada a superheroapi"}))
	}
	//Adiciona fechamento da resposta a pilha do defer
	defer resp.Body.Close()
	//Lê resposta em uma array de bytes
	respByte, err := ioutil.ReadAll(resp.Body)
	//Tratamento de erro da leitura de resposta
	if err != nil {
		return c.Render(http.StatusOK, r.JSON(map[string]string{"message": "Erro na leitura da resposta"}))
	}
	//Cria variavel para json de resposta e deserializa respByte
	var searchResponse SearchResponse
	err = json.Unmarshal(respByte, &searchResponse)
	//Tratamento de erro da deserialização
	if err != nil {
		return c.Render(http.StatusOK, r.JSON(map[string]string{"message": "Erro na conversão de Json"}))
	}
	if searchResponse.Response == "error" {
		return c.Render(http.StatusOK, r.JSON(searchResponse))
	}
	//Pega conexão ao banco de dados do contexto
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	//Para cada resultado da pesquisa, adiciona id da superheroapi em query
	results := searchResponse.Results
	var resultsIDS []interface{}
	for i, result := range results {
		//Converte id do resultado de string para int
		originalID, err := strconv.Atoi(result.ID)
		//Tratamento de erro da conversão
		if err != nil {
			message := fmt.Sprintf("Erro na conversão do resultado %d", i)
			return c.Render(http.StatusOK, r.JSON(map[string]string{"message": message}))
		}
		//Adiciona id a slice
		resultsIDS = append(resultsIDS, originalID)
	}
	q := tx.Where("original_id in (?)", resultsIDS...)
	//Seleciona apenas a coluna com as ids da superheroapi
	q.Select("original_id")
	//Array para conter as ids da superheroapi que ja estão no banco
	//var alreadyOnDB []int
	//Executa query para pegar quais ids ja estão no banco de dados
	supersAlreadyOnDB := []models.Super{}
	q.All(&supersAlreadyOnDB)
	//Array para devolver supers registrados
	var registredSupers []models.Super = []models.Super{}
	//Para cada resultado da pesquisa, confere se o id da superhero api bate com um original_id já no banco de dados, e adiciona novo super ao banco de dados caso não
	for i, result := range results {
		//Boolean que indica se o id da superapi do resultado atual já esta no bando da dados
		idAlreadyOnDB := false
		//Converte id da superapi do resultado atual de string para int
		currentID, err := strconv.Atoi(result.ID)
		//Tratamento de erro da conversão
		if err != nil {
			message := fmt.Sprintf("Erro na conversão do resultado %d", i)
			return c.Render(http.StatusOK, r.JSON(map[string]string{"message": message}))
		}
		//Para cada id da superapi achado no bando de dados, confere se o resultado atual possui o mesmo id
		for _, superOnDB := range supersAlreadyOnDB {
			if currentID == superOnDB.OriginalID {
				idAlreadyOnDB = true
				break
			}
		}
		//Caso o resultado atual não ja esteja no banco de dados, cria novo super
		if !idAlreadyOnDB {
			//Atribui valores do resultado ao objeto de modelo de super
			super := &models.Super{}
			super.OriginalID = currentID
			super.Name = result.Name
			if !isNullValue(result.Biography.FullName) {
				super.FullName = result.Biography.FullName
			}
			if !isNullValue(result.Biography.PlaceOfBirth) {
				super.PlaceOfBirth = result.Biography.PlaceOfBirth
			}
			if !isNullValue(result.Biography.FirstAppearance) {
				super.FirstAppearance = result.Biography.FirstAppearance
			}
			super.AlterEgos = result.Biography.AlterEgos
			super.Publisher = result.Biography.Publisher
			super.Alignment = result.Biography.Alignment
			if !isNullValue(result.Appearance.Gender) {
				super.Gender = result.Appearance.Gender
			}
			if !isNullValue(result.Appearance.Race) {
				super.Race = result.Appearance.Race
			}
			super.HeightFeet = result.Appearance.Height[0]
			if len(result.Appearance.Height) > 1 {
				heightCm, err := convertFirstWordToInt(result.Appearance.Height[1])
				if err != nil {
					super.HeightCm = 0
				} else {
					super.HeightCm = heightCm
				}
			}
			super.WeightLb = result.Appearance.Weight[0]
			if len(result.Appearance.Weight) > 1 {
				weightKg, err := convertFirstWordToInt(result.Appearance.Weight[1])
				if err != nil {
					super.WeightKg = 0
				} else {
					super.WeightKg = weightKg
				}
			}
			if !isNullValue(result.Appearance.EyeColor) {
				super.EyeColor = result.Appearance.EyeColor
			}
			if !isNullValue(result.Appearance.HairColor) {
				super.HairColor = result.Appearance.HairColor
			}
			if !isNullValue(result.Work.Occupation) {
				super.Occupation = result.Work.Occupation
			}
			if !isNullValue(result.Work.Base) {
				super.Base = result.Work.Base
			}
			if !isNullValue(result.Image.URL) {
				super.Image = result.Image.URL
			}
			intelligence, err := convertFirstWordToInt(result.Powerstats.Intelligence)
			if err != nil {
				super.Intelligence = 0
			} else {
				super.Intelligence = intelligence
			}
			strength, err := convertFirstWordToInt(result.Powerstats.Strength)
			if err != nil {
				super.Strength = 0
			} else {
				super.Strength = strength
			}
			speed, err := convertFirstWordToInt(result.Powerstats.Speed)
			if err != nil {
				super.Speed = 0
			} else {
				super.Speed = speed
			}
			durability, err := convertFirstWordToInt(result.Powerstats.Durability)
			if err != nil {
				super.Durability = 0
			} else {
				super.Durability = durability
			}
			power, err := convertFirstWordToInt(result.Powerstats.Power)
			if err != nil {
				super.Power = 0
			} else {
				super.Power = power
			}
			combat, err := convertFirstWordToInt(result.Powerstats.Combat)
			if err != nil {
				super.Combat = 0
			} else {
				super.Combat = combat
			}
			//Valida e cria super
			tx.ValidateAndCreate(super)
			registredSupers = append(registredSupers, *super)
		}
	}
	return c.Render(http.StatusOK, r.JSON(registredSupers))
}

//SupersAll lista todos os supers registrados
func SupersAll(c buffalo.Context) error {
	//Pega conexão ao banco de dados do contexto
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	//Cria variavel para receber supers
	supers := &models.Supers{}
	//Executa query para pegar todos os supers
	tx.All(supers)
	//Renderiza json
	return c.Render(http.StatusOK, r.JSON(supers))
}

//SupersHeros lista todos os supers registrados que são herois
func SupersHeros(c buffalo.Context) error {
	//Pega conexão ao banco de dados do contexto
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	//Cria array de supers
	supers := &models.Supers{}
	//Cria query que seleciona apenas supers que são herois
	q := tx.Where("alignment = ?", "good")
	//Executa query
	q.All(supers)
	//Renderiza json
	return c.Render(http.StatusOK, r.JSON(supers))
}

//SupersVillains lista todos os supers registrados que são vilões
func SupersVillains(c buffalo.Context) error {
	//Pega conexão ao banco de dados do contexto
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	//Cria array de supers
	supers := &models.Supers{}
	//Cria query que seleciona apenas supers que são vilões
	q := tx.Where("alignment = ?", "bad")
	//Executa query
	q.All(supers)
	//Renderiza json
	return c.Render(http.StatusOK, r.JSON(supers))
}

//SupersDestroy deleta do banco de dados o super cujo id é passado pelo parametro "super_id"
func SupersDestroy(c buffalo.Context) error {
	//Pega conexão ao banco de dados do contexto
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	//Aloca um super vazio
	super := &models.Super{}

	//Procura super
	if err := tx.Find(super, c.Param("super_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(super); err != nil {
		return err
	}
	return c.Render(http.StatusOK, r.JSON(super))
}

//SupersSearch é a ação de busca de supers
func SupersSearch(c buffalo.Context) error {
	//Pega conexão ao banco de dados do contexto
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	//Query de busca
	q := tx.Q()
	//Parametros enviados pela rota
	params := c.Params()
	//Adicionar where para cada parametro
	//Confere parametros de powerstats
	intelligenceParam := params.Get("intelligence")
	if intelligenceParam != "" {
		intelligence, err := convertFirstWordToInt(intelligenceParam)
		if err != nil {
			return c.Render(http.StatusBadRequest, r.JSON(map[string]string{"Erro": "Parametro intelligence deve conter int"}))
		}
		q.Where("intelligence >= ?", intelligence)
	}
	strengthParam := params.Get("strength")
	if strengthParam != "" {
		strength, err := convertFirstWordToInt(strengthParam)
		if err != nil {
			return c.Render(http.StatusBadRequest, r.JSON(map[string]string{"Erro": "Parametro strength deve conter int"}))
		}
		q.Where("strength >= ?", strength)
	}
	speedParam := params.Get("speed")
	if speedParam != "" {
		speed, err := convertFirstWordToInt(speedParam)
		if err != nil {
			return c.Render(http.StatusBadRequest, r.JSON(map[string]string{"Erro": "Parametro speed deve conter int"}))
		}
		q.Where("speed >= ?", speed)
	}
	durabilityParam := params.Get("durability")
	if durabilityParam != "" {
		durability, err := convertFirstWordToInt(durabilityParam)
		if err != nil {
			return c.Render(http.StatusBadRequest, r.JSON(map[string]string{"Erro": "Parametro durability deve conter int"}))
		}
		q.Where("durability >= ?", durability)
	}
	powerParam := params.Get("power")
	if powerParam != "" {
		power, err := convertFirstWordToInt(powerParam)
		if err != nil {
			return c.Render(http.StatusBadRequest, r.JSON(map[string]string{"Erro": "Parametro power deve conter int"}))
		}
		q.Where("power >= ?", power)
	}
	combatParam := params.Get("combat")
	if combatParam != "" {
		combat, err := convertFirstWordToInt(combatParam)
		if err != nil {
			return c.Render(http.StatusBadRequest, r.JSON(map[string]string{"Erro": "Parametro combat deve conter int"}))
		}
		q.Where("combat >= ?", combat)
	}
	//Aloca array de supers para receber resultado da query
	supers := []models.Super{}
	//Executa query
	q.All(&supers)
	//Confere parametros para filtrar resultado da query
	//Carrega parametro name
	name := params.Get("name")
	if name != "" {
		supers = filterSupersName(supers, name)
	}
	//Carrega parametro full_name
	fullName := params.Get("full_name")
	if fullName != "" {
		supers = filterSupersFullName(supers, fullName)
	}
	//Carrega parametro occupation
	occupation := params.Get("occupation")
	if occupation != "" {
		supers = filterSupersOccupation(supers, occupation)
	}
	return c.Render(http.StatusOK, r.JSON(supers))
}
