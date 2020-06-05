package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"super_api/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
)

const relativeNameRegex = `[^,\(\)]*`
const relativeRelationRegex = `(\([^\)]*\))?`
const relativeRegex = relativeNameRegex + relativeRelationRegex

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

func filterSupersGender(supers []models.Super, name string) []models.Super {
	//Aloca array de supers
	filteredSupers := []models.Super{}
	//Para cada super, comparar name
	for _, super := range supers {
		if strings.ToLower(super.Gender) == strings.ToLower(name) {
			filteredSupers = append(filteredSupers, super)
		}
	}
	return filteredSupers
}

func filterSupersRace(supers []models.Super, name string) []models.Super {
	//Aloca array de supers
	filteredSupers := []models.Super{}
	//Para cada super, comparar name
	for _, super := range supers {
		if stringContainsSubstring(super.Race, name) {
			filteredSupers = append(filteredSupers, super)
		}
	}
	return filteredSupers
}

func filterSupersEyeColor(supers []models.Super, name string) []models.Super {
	//Aloca array de supers
	filteredSupers := []models.Super{}
	//Para cada super, comparar name
	for _, super := range supers {
		if stringContainsSubstring(super.EyeColor, name) {
			filteredSupers = append(filteredSupers, super)
		}
	}
	return filteredSupers
}

func filterSupersHairColor(supers []models.Super, name string) []models.Super {
	//Aloca array de supers
	filteredSupers := []models.Super{}
	//Para cada super, comparar name
	for _, super := range supers {
		if stringContainsSubstring(super.HairColor, name) {
			filteredSupers = append(filteredSupers, super)
		}
	}
	return filteredSupers
}

func filterSupersAlias(supers []models.Super, name string) []models.Super {
	//Aloca array de supers
	filteredSupers := []models.Super{}
	//Para cada super, comparar name
	for _, super := range supers {
		for _, alias := range super.Aliases {
			if stringContainsSubstring(alias.Name, name) {
				filteredSupers = append(filteredSupers, super)
				break
			}
		}
	}
	return filteredSupers
}

func filterSupersGroup(supers []models.Super, name string) []models.Super {
	//Aloca array de supers
	filteredSupers := []models.Super{}
	//Para cada super, comparar name
	for _, super := range supers {
		for _, group := range super.Groups {
			if stringContainsSubstring(group.Name, name) {
				filteredSupers = append(filteredSupers, super)
				break
			}
		}
	}
	return filteredSupers
}

func filterSupersRelative(supers []models.Super, name string) []models.Super {
	//Aloca array de supers
	filteredSupers := []models.Super{}
	//Para cada super, comparar name
	for _, super := range supers {
		for _, relative := range super.Relatives {
			if stringContainsSubstring(relative.Name, name) {
				filteredSupers = append(filteredSupers, super)
				break
			}
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
		return c.Render(http.StatusBadRequest, r.JSON(map[string]string{"message": "Sem parametro name para buscar"}))
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
			//Caso o resultado tenha aliases
			if len(result.Biography.Aliases) > 0 {
				//Aloca slice vazio de aliases
				super.Aliases = []models.Alias{}
				//Para cada alias do resultado
				for _, aliasName := range result.Biography.Aliases {
					//Se não for um valor invalido
					if !isNullValue(aliasName) {
						//Aloca alias vazio
						alias := models.Alias{}
						//Confere se alias ja existe
						tx.Where("name = ?", aliasName).First(&alias)
						//Se existir, adicionar ao super atual, senão criar novo alias e adicionar ao super atual
						if alias.Name != "" {
							super.Aliases = append(super.Aliases, alias)
						} else {
							alias.Name = aliasName
							super.Aliases = append(super.Aliases, alias)
						}
					}
				}
			}
			//Groups
			//Separa string de groups
			firstSplits := strings.Split(result.Connections.GroupAffiliation, ",")
			//Aloca slice vazio de strings
			groups := []string{}
			//Para cada slice em firstSplists, dividir em ";" e adicionar a groups
			for _, firstSplit := range firstSplits {
				splits := strings.Split(firstSplit, ";")
				groups = append(groups, splits...)
			}
			//Para cada grupo
			for _, groupName := range groups {
				//Se o group for valido
				if !isNullValue(groupName) {
					//Retira espaços do inicio e fim
					groupName := strings.TrimSpace(groupName)
					//Aloca grupo vazio
					group := models.Group{}
					//Confere se o grupo ja existe
					tx.Where("name = ?", groupName).First(&group)
					//Se existir, adicionar ao super atual, senão criar novo grupo e adicionar ao super atual
					if group.Name != "" {
						super.Groups = append(super.Groups, group)
					} else {
						group.Name = groupName
						super.Groups = append(super.Groups, group)
					}
				}
			}
			//Relatives
			//Se possuir valor valido
			if !isNullValue(result.Connections.Relatives) {
				//Aloca slice vazio de relatives
				super.Relatives = []models.Relative{}
				re := regexp.MustCompile(relativeRegex)
				relativesStrings := re.FindAllString(result.Connections.Relatives, -1)
				//Para cada relative achado na expressão regular
				for _, relativeString := range relativesStrings {
					//Se o relative for valido
					if !isNullValue(relativeString) {
						//Aloca relative vazio
						relative := models.Relative{}
						//Retira espaços do inicio e fim
						relativeString := strings.TrimSpace(relativeString)
						//Confere se o relative ja existe
						tx.Where("name = ?", relativeString).First(&relative)
						//Se existir, adicionar ao super atual, senão criar antes de adicionar
						if relative.Name != "" {
							super.Relatives = append(super.Relatives, relative)
						} else {
							relative.Name = relativeString
							super.Relatives = append(super.Relatives, relative)
						}
					}
				}
			}
			//Valida e cria super
			tx.Eager().ValidateAndCreate(super)
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
	tx.Eager().All(supers)
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
	q.Eager().All(supers)
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
	q.Eager().All(supers)
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
	//Parametros enviados pela rota
	params := c.Params()
	//Confere se existe parametro uuid
	uuid := params.Get("uuid")
	if uuid != "" {
		super := models.Super{}
		tx.Find(&super, uuid)
		return c.Render(http.StatusOK, r.JSON(super))
	}
	//Query de busca
	q := tx.Q()
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
	//Confere parametros de altura e peso
	//Altura minima
	mimHeightParam := params.Get("min_height")
	if mimHeightParam != "" {
		minHeight, err := convertFirstWordToInt(mimHeightParam)
		if err != nil {
			return c.Render(http.StatusBadRequest, r.JSON(map[string]string{"Erro": "Parametro mim_height deve conter int"}))
		}
		q.Where("height_cm >= ?", minHeight)
	}
	//Altura maxima
	maxHeightParam := params.Get("max_height")
	if maxHeightParam != "" {
		maxHeight, err := convertFirstWordToInt(maxHeightParam)
		if err != nil {
			return c.Render(http.StatusBadRequest, r.JSON(map[string]string{"Erro": "Parametro max_height deve conter int"}))
		}
		q.Where("height_cm <= ?", maxHeight)
	}
	//Peso minimo
	mimWeightParam := params.Get("min_weight")
	if mimWeightParam != "" {
		minWeight, err := convertFirstWordToInt(mimWeightParam)
		if err != nil {
			return c.Render(http.StatusBadRequest, r.JSON(map[string]string{"Erro": "Parametro mim_weight deve conter int"}))
		}
		q.Where("weight_kg >= ?", minWeight)
	}
	//Peso maximo
	maxWeightParam := params.Get("max_weight")
	if maxWeightParam != "" {
		maxWeight, err := convertFirstWordToInt(maxWeightParam)
		if err != nil {
			return c.Render(http.StatusBadRequest, r.JSON(map[string]string{"Erro": "Parametro max_weight deve conter int"}))
		}
		q.Where("weight_kg <= ?", maxWeight)
	}
	//Imagem
	image := params.Get("image")
	if image != "" {
		q.Where("image = ?", image)
	}
	//Aloca array de supers para receber resultado da query
	supers := []models.Super{}
	//Executa query
	q.Eager().All(&supers)
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
	//Carrega parametro alias
	alias := params.Get("alias")
	if alias != "" {
		supers = filterSupersAlias(supers, alias)
	}
	//Carrega parametro group
	group := params.Get("group")
	if group != "" {
		supers = filterSupersGroup(supers, group)
	}
	//Carega parametro relative
	relative := params.Get("relative")
	if relative != "" {
		supers = filterSupersRelative(supers, relative)
	}
	//Carrega parametro occupation
	occupation := params.Get("occupation")
	if occupation != "" {
		supers = filterSupersOccupation(supers, occupation)
	}
	//Carrega parametro gender
	gender := params.Get("gender")
	if gender != "" {
		supers = filterSupersGender(supers, gender)
	}
	//Carrega parametro race
	race := params.Get("race")
	if race != "" {
		supers = filterSupersRace(supers, race)
	}
	//Carrega parametro eye_color
	eyeColor := params.Get("eye_color")
	if eyeColor != "" {
		supers = filterSupersEyeColor(supers, eyeColor)
	}
	//Carrega parametro hair_color
	hairColor := params.Get("hair_color")
	if hairColor != "" {
		supers = filterSupersHairColor(supers, hairColor)
	}
	return c.Render(http.StatusOK, r.JSON(supers))
}
