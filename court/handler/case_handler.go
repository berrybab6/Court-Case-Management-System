package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	entity "github.com/Surafeljava/Court-Case-Management-System/Entity"
	"github.com/Surafeljava/Court-Case-Management-System/caseUse"
	"github.com/Surafeljava/Court-Case-Management-System/form"
)

type CaseHandler struct {
	tmpl    *template.Template
	caseSrv caseUse.CaseService
}

func NewCaseHandler(T *template.Template, CS caseUse.CaseService) *CaseHandler {
	return &CaseHandler{tmpl: T, caseSrv: CS}
}

//Get all the cases in the court to the admin
func (lh *CaseHandler) Cases(w http.ResponseWriter, r *http.Request) {
	cases := []entity.Case{}
	cases, err := lh.caseSrv.Cases()
	if err != nil {
		panic(err)
	}
	lh.tmpl.ExecuteTemplate(w, "admin.cases.layout", cases)
}

//Add new Cases
func (lh *CaseHandler) NewCase(w http.ResponseWriter, r *http.Request) {
	//lh.tmpl.ExecuteTemplate(w, "admin.newcase.layout", nil)
	//TODO Create a new case here...

	if r.Method == http.MethodPost {
		//TODO get the number of cases and generate a new case number
		//Get the case number
		allCases, err1 := lh.caseSrv.Cases()
		if err1 != nil {
			panic(err1)
		}

		cs_num := len(allCases) + 1
		case_num := fmt.Sprintf("CS%d", cs_num)

		case_title := r.FormValue("case_title")
		case_desc := r.FormValue("case_desc")
		case_type := r.FormValue("case_type")
		court_date := r.FormValue("court_date")
		case_judge := r.FormValue("case_judge")

		//validate the inputs from the form

		newCaseForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		newCaseForm.Required("case_title", "case_desc", "case_type", "court_date", "case_judge")
		newCaseForm.MinLength("case_desc", 8)
		newCaseForm.MinLength("case_judge", 3)
		// newCaseForm.CSRF = token

		//Checking the validation of the form inputs
		if !newCaseForm.Valid() {
			lh.tmpl.ExecuteTemplate(w, "login.layout", newCaseForm)
			return
		}

		the_court_date, _ := time.Parse("2006-01-02", court_date)
		//the_case_creation, _ := time.Parse("2006-01-02", time.Now())

		newcs := entity.Case{CaseNum: case_num, CaseTitle: case_title, CaseDesc: case_desc, CaseStatus: "open", CaseType: case_type, CaseCreation: time.Now(), CaseCourtDate: the_court_date, CaseJudge: case_judge}

		err2 := lh.caseSrv.CreateCase(&newcs)

		if len(err2) > 0 {
			panic(err2)
		}

	} else if r.Method == http.MethodGet {
		case_type := r.URL.Query().Get("case_type")

		jud, _ := lh.caseSrv.CaseJudges(case_type)
		lh.tmpl.ExecuteTemplate(w, "admin.newcase.layout", jud)
		//lh.tmpl.ExecuteTemplate(w, "admin.newcase.layout", nil)
	} else {
		lh.tmpl.ExecuteTemplate(w, "admin.newcase.layout", nil)
	}

	http.Redirect(w, r, "/admin/cases", http.StatusSeeOther)
}

//Delete existing cases
func (lh *CaseHandler) DeleteCase(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")
		//idr := 30
		id, err := strconv.Atoi(idRaw)

		if err != nil {
			fmt.Println("Can't get the id!")
			panic(err)
		}

		errs := lh.caseSrv.DeleteCase(id)

		if len(errs) > 0 {
			http.Redirect(w, r, "/admin/cases", http.StatusSeeOther)
		}

		http.Redirect(w, r, "/admin/cases", http.StatusSeeOther)

	}
}

//Update existing cases
func (lh *CaseHandler) UpdateCase(w http.ResponseWriter, r *http.Request) {
	//lh.tmpl.ExecuteTemplate(w, "admin.newcase.layout", nil)
	//TODO Update case here...
	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")
		//idr := 30
		id, err := strconv.Atoi(idRaw)

		if err != nil {
			fmt.Println("Can't get the id!")
			panic(err)
		}

		cs, _ := lh.caseSrv.Case(id)

		lh.tmpl.ExecuteTemplate(w, "admin.case.update.layout", cs)

	} else if r.Method == http.MethodPost {

		id, _ := strconv.Atoi(r.FormValue("id"))
		case_num := r.FormValue("case_num")
		case_title := r.FormValue("case_title")
		case_desc := r.FormValue("case_desc")
		case_type := r.FormValue("case_type")
		case_creation := time.Now()
		case_court_date, _ := time.Parse("2006-01-02", r.FormValue("court_date"))
		case_judge := r.FormValue("case_judge")

		cs := entity.Case{ID: uint(id), CaseNum: case_num, CaseTitle: case_title, CaseDesc: case_desc, CaseType: case_type, CaseCreation: case_creation, CaseCourtDate: case_court_date, CaseJudge: case_judge}
		_, err := lh.caseSrv.UpdateCase(&cs)

		if len(err) > 0 {
			fmt.Println(">> -- >> -- >> error on sending the Case to the updateCase func")
			panic(err)
		}

		http.Redirect(w, r, "/admin/cases", http.StatusSeeOther)

	} else {

		http.Redirect(w, r, "/admin/cases", http.StatusSeeOther)
	}
}

func (lh *CaseHandler) CaseTypeJudge(w http.ResponseWriter, r *http.Request) {
	//Get the judges suitable for that case
}

//Close a case by adding final decision and description >> by the judge
func (lh *CaseHandler) CloseCase(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		juid := r.URL.Query().Get("id")

		cs, er := lh.caseSrv.JudgeCases(juid)

		if er != nil {
			fmt.Println("Error Getting Cases for the judge")
		}

		lh.tmpl.ExecuteTemplate(w, "judge.case.close.layout", cs)

	} else if r.Method == http.MethodPost {

		case_num := r.FormValue("case_num")

		decision_date, _ := time.Parse("2006-02-06", time.Now().String())
		decision := r.FormValue("final_decision")
		decision_desc := r.FormValue("close_desc")

		//cs := entity.Case{ID: uint(id), CaseNum: case_num, CaseTitle: case_title, CaseDesc: case_desc, CaseType: case_type, CaseCreation: case_creation, CaseCourtDate: case_court_date, CaseJudge: case_judge}
		des := entity.Decision{CaseNum: case_num, DecisionDate: decision_date, Decision: decision, DecisionDesc: decision_desc}
		err := lh.caseSrv.CloseCase(case_num, &des)

		if len(err) > 0 {
			fmt.Println(">> -- >> -- >> error on sending the Case to the CloseCase func")
			panic(err)
		}

		http.Redirect(w, r, "/judge/cases/close", http.StatusSeeOther)
		//lh.tmpl.ExecuteTemplate(w, "judge.home.layout", nil)

	} else {

		http.Redirect(w, r, "/judge/cases/close", http.StatusSeeOther)
	}
}
