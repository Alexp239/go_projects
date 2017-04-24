package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

type XMLUser struct {
	Id        int    `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

type Query struct {
	XMLName xml.Name  `xml:"root"`
	Users   []XMLUser `xml:"row"`
}

type IdAscUserSorter []User

func (a IdAscUserSorter) Len() int           { return len(a) }
func (a IdAscUserSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a IdAscUserSorter) Less(i, j int) bool { return a[i].Id < a[j].Id }

type IdDecUserSorter []User

func (a IdDecUserSorter) Len() int           { return len(a) }
func (a IdDecUserSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a IdDecUserSorter) Less(i, j int) bool { return a[i].Id > a[j].Id }

type AgeAscUserSorter []User

func (a AgeAscUserSorter) Len() int           { return len(a) }
func (a AgeAscUserSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a AgeAscUserSorter) Less(i, j int) bool { return a[i].Age < a[j].Age }

type AgeDecUserSorter []User

func (a AgeDecUserSorter) Len() int           { return len(a) }
func (a AgeDecUserSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a AgeDecUserSorter) Less(i, j int) bool { return a[i].Age > a[j].Age }

type NameAscUserSorter []User

func (a NameAscUserSorter) Len() int           { return len(a) }
func (a NameAscUserSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a NameAscUserSorter) Less(i, j int) bool { return a[i].Name < a[j].Name }

type NameDecUserSorter []User

func (a NameDecUserSorter) Len() int           { return len(a) }
func (a NameDecUserSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a NameDecUserSorter) Less(i, j int) bool { return a[i].Name > a[j].Name }

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getIDFromUsers(a []User) (res []int) {
	for _, u := range a {
		res = append(res, u.Id)
	}
	return
}

func getNameFromUsers(a []User) (res []string) {
	for _, u := range a {
		res = append(res, u.Name)
	}
	return
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	xmlFile, err := os.Open("dataset.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	data, _ := ioutil.ReadAll(xmlFile)

	var q Query
	xml.Unmarshal(data, &q)

	searcherParams, _ := url.ParseQuery(r.URL.String()[2:])
	limit, err := strconv.Atoi(searcherParams["limit"][0])
	PanicOnErr(err)
	offset, err := strconv.Atoi(searcherParams["offset"][0])
	PanicOnErr(err)
	query := searcherParams["query"][0]
	orderField := searcherParams["order_field"][0]
	orderBy := searcherParams["order_by"][0]

	var tmp, res []User

	for _, xmlUser := range q.Users {
		user := User{
			Id:     xmlUser.Id,
			Name:   xmlUser.LastName + " " + xmlUser.FirstName,
			About:  xmlUser.About,
			Age:    xmlUser.Age,
			Gender: xmlUser.Gender,
		}
		tmp = append(tmp, user)
	}

	for _, user := range tmp {
		if strings.Contains(user.About, query) || strings.Contains(user.Name, query) {
			res = append(res, user)
		}
	}

	switch orderField {
	case "Id":
		if orderBy == "0" {
			sort.Sort(IdDecUserSorter(res))
		} else {
			sort.Sort(IdAscUserSorter(res))
		}
	case "Age":
		if orderBy == "0" {
			sort.Sort(AgeDecUserSorter(res))
		} else {
			sort.Sort(AgeAscUserSorter(res))
		}
	case "Name":
		if orderBy == "0" {
			sort.Sort(NameDecUserSorter(res))
		} else {
			sort.Sort(NameAscUserSorter(res))
		}
	default:
		sort.Sort(NameAscUserSorter(res))
	}

	len := len(res)

	data, _ = json.Marshal(res[min(len, offset):min(len, offset+limit)])

	w.Write(data)
}

func TestSearch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	expected := SearchResponse{Users: []User{User{
		Id:     32,
		Name:   "Knapp Christy",
		Age:    40,
		About:  "Incididunt culpa dolore laborum cupidatat consequat. Aliquip cupidatat pariatur sit consectetur laboris labore anim labore. Est sint ut ipsum dolor ipsum nisi tempor in tempor aliqua. Aliquip labore cillum est consequat anim officia non reprehenderit ex duis elit. Amet aliqua eu ad velit incididunt ad ut magna. Culpa dolore qui anim consequat commodo aute.\n",
		Gender: "female",
	}, User{
		Id:     13,
		Name:   "Davidson Whitley",
		Age:    40,
		About:  "Consectetur dolore anim veniam aliqua deserunt officia eu. Et ullamco commodo ad officia duis ex incididunt proident consequat nostrud proident quis tempor. Sunt magna ad excepteur eu sint aliqua eiusmod deserunt proident. Do labore est dolore voluptate ullamco est dolore excepteur magna duis quis. Quis laborum deserunt ipsum velit occaecat est laborum enim aute. Officia dolore sit voluptate quis mollit veniam. Laborum nisi ullamco nisi sit nulla cillum et id nisi.\n",
		Gender: "male",
	}},
		NextPage: true,
	}
	res, _ := doSearch(ts.URL, 2, 0, "", "Age", 0)
	if !reflect.DeepEqual(*res, expected) {
		t.Errorf("expected \n%+v\n, got \n%+v\n", expected, *res)
	}

	res, err := doSearch(ts.URL, -1, 0, "", "Age", 0)
	if err != nil && err.Error() != errors.New("limit must be > 0").Error() {
		t.Errorf("expected \n%+v\n, got \n%+v\n", errors.New("limit must be > 0").Error(), err.Error())
	}

	res, err = doSearch(ts.URL, 1, -1, "", "Age", 0)
	if err != nil && err.Error() != errors.New("offset must be > 0").Error() {
		t.Errorf("expected \n%+v\n, got \n%+v\n", errors.New("offset must be > 0").Error(), err.Error())
	}

	expected = SearchResponse{Users: []User{User{
		Id:     31,
		Name:   "Scott Palmer",
		Age:    37,
		About:  "Elit fugiat commodo laborum quis eu consequat. In velit magna sit fugiat non proident ipsum tempor eu. Consectetur exercitation labore eiusmod occaecat adipisicing irure consequat fugiat ullamco aliquip nostrud anim irure enim. Duis do amet cillum eiusmod eu sunt. Minim minim sunt sit sit enim velit sint tempor enim sint aliquip voluptate reprehenderit officia. Voluptate magna sit consequat adipisicing ut eu qui.\n",
		Gender: "male",
	},
	},
		NextPage: false,
	}
	res, err = doSearch(ts.URL, 200, 0, "velit magna sit fugiat", "", 0)
	if err != nil || !reflect.DeepEqual(*res, expected) {
		t.Errorf("expected \n%+v\n, got \n%+v\n", expected, *res)
	}

	res, err = doSearch("vk.com/", 200, 0, "velit magna sit fugiat", "", 0)
	err, ok := err.(net.Error)
	if !ok {
		t.Errorf("expected net.Error, got \n%+v\n", err)
	}

	expectedID := []int{1, 2}
	res, err = doSearch(ts.URL, 2, 1, "", "Id", 1)
	if err != nil || res.NextPage != true || !reflect.DeepEqual(getIDFromUsers(res.Users), expectedID) {
		t.Errorf("expected \n%+v\n, got \n%+v\n", expected, *res)
	}

	expectedNames := []string{"Aguilar Brooks", "Anderson Gonzalez"}
	res, err = doSearch(ts.URL, 2, 0, "", "Name", 1)
	if err != nil || res.NextPage != true || !reflect.DeepEqual(getNameFromUsers(res.Users), expectedNames) {
		t.Errorf("expected \n%+v\n, got \n%+v\n", expected, *res)
	}
}

func SearchServerTest(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(XMLUser{})

	w.Write(data)
}

func SearchServerTestTimeout(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 2)
}

func TestServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServerTest))
	// expectedErr := errors.New("json: cannot unmarshal object into Go value of type []main.User")
	res, err := doSearch(ts.URL, 2, 1, "", "Id", 1)
	expected := json.UnmarshalTypeError{Value: "object", Type: reflect.TypeOf([]User{})}
	if res != nil || err.Error() != expected.Error() {
		t.Errorf("expected err:\n%+v\n, got err:\n%+v\n", expected.Error(), err.Error())
	}
	ts.Close()
	ts = httptest.NewServer(http.HandlerFunc(SearchServerTestTimeout))
	res, err = doSearch(ts.URL, 2, 1, "", "Id", 1)
	if res != nil {
		if err, ok := err.(net.Error); !ok || err.Timeout() {
			t.Errorf("expected err:\n%+v\n, got err:\n%+v\n", expected.Error(), err.Error())
		}
	}
	ts.Close()
}
