// methods.simulation.go
// class methods of the objects specified in models.simulation.go
package models

import (
	"fmt"
	"log"
	"strconv"
)

//METHODS OF INDUSTRIES

// crude searches without database implementation
// justified because we can avoid the complications of a database implementation
// and the size of the tables is not large, because they are provided on a per-user basis
// However as the simulations get large, this may become more problematic (let's find out pragmatically)
// In that case some more sophisticated system, such as a local database, may be needed
// A simple solution would be to add direct links to related objects in the models
// perhaps populated by an asynchronous process in the background

// A default Industry_stock returned if any condition is not met (that is, if the predicated stock does not exist)
// Used to signal to the user that there has been a programme error
var NotFoundIndustryStock = Industry_Stock{
	Id:            0,
	Simulation_id: 0,
	Commodity_id:  0,
	Name:          "NOT FOUND",
	Usage_type:    "PROGRAMME ERROR",
	Size:          -1,
	Value:         -1,
	Price:         -1,
	Requirement:   -1,
	Demand:        -1,
}

// A default Industry_stock returned if any condition is not met (that is, if the predicated stock does not exist)
// Used to signal to the user that there has been a programme error
var NotFoundClassStock = Class_Stock{
	Id:            0,
	Simulation_id: 0,
	Commodity_id:  0,
	Name:          "NOT FOUND",
	Usage_type:    "PROGRAMME ERROR",
	Size:          -1,
	Value:         -1,
	Price:         -1,
	Demand:        -1,
}

var NotFoundCommodity = Commodity{
	Id:                          0,
	Name:                        "NOT FOUND",
	Simulation_id:               0,
	Time_Stamp:                  0,
	Origin:                      "UNDEFINED",
	Usage:                       "UNDEFINED",
	Size:                        0,
	Total_Value:                 0,
	Total_Price:                 0,
	Unit_Value:                  0,
	Unit_Price:                  0,
	Turnover_Time:               0,
	Demand:                      0,
	Supply:                      0,
	Allocation_Ratio:            0,
	Display_Order:               0,
	Image_Name:                  "UNDEFINED",
	Tooltip:                     "UNDEFINED",
	Monetarily_Effective_Demand: 0,
	Investment_Proportion:       0,
}

func (p Pair) Display() string {
	return fmt.Sprintf("%0.2f", p.Viewed)
}

// returns the money stock of the given industry
func (industry Industry) MoneyStock() Industry_Stock {
	username := industry.UserName
	stockList := *Users[username].IndustryStocks()
	for i := 0; i < len(stockList); i++ {
		s := stockList[i]
		if (s.Industry_id == industry.Id) && (s.Usage_type == `Money`) {
			return s
		}
	}
	return NotFoundIndustryStock
}

// returns the sales stock of the given industry
func (industry Industry) SalesStock() Industry_Stock {
	username := industry.UserName
	stockList := *Users[username].IndustryStocks()
	for i := 0; i < len(stockList); i++ {
		s := &stockList[i]
		if (s.Industry_id == industry.Id) && (s.Usage_type == `Sales`) {
			return *s
		}
	}
	return NotFoundIndustryStock
}

// returns the Labour Power stock of the given industry
// bit of a botch to use the name of the commodity as a search term
func (industry Industry) VariableCapital() Industry_Stock {
	username := industry.UserName
	stockList := *Users[username].IndustryStocks()
	for i := 0; i < len(stockList); i++ {
		s := &stockList[i]
		if (s.Industry_id == industry.Id) && (s.Usage_type == `Production`) && (s.CommodityName() == "Labour Power") {
			return *s
		}
	}
	return NotFoundIndustryStock
}

// returns the commodity that an industry produces
func (industry Industry) OutputCommodity() *Commodity {
	return industry.SalesStock().Commodity()
}

// return the productive capital stock of the given industry
// under development - at present assumes there is only one
func (industry Industry) ConstantCapital() Industry_Stock {
	username := industry.UserName
	stockList := *Users[username].IndustryStocks()
	for i := 0; i < len(stockList); i++ {
		s := &stockList[i]
		if (s.Industry_id == industry.Id) && (s.Usage_type == `Production`) && (s.CommodityName() == "Means of Production") {
			return *s
		}
	}
	return NotFoundIndustryStock
}

// returns all the constant capitals of a given industry.
// Under development.
// func (industry Industry) ConstantCapitals() []Stock {
// 	return &stocks [Programming error here]
// }

// METHODS OF SOCIAL CLASSES

// returns the sales stock of the given class
// was 	err = db.SDB.QueryRowx("SELECT * FROM stocks where Owner_Id = ? AND Usage_type =?", class.Id, "Sales").StructScan(&stock)
func (class Class) MoneyStock() Class_Stock {
	username := class.UserName
	stockList := *Users[username].ClassStocks()

	for i := 0; i < len(stockList); i++ {
		s := &stockList[i]
		if (s.Class_id == class.Id) && (s.Usage_type == `Money`) {
			return *s
		}
	}
	return NotFoundClassStock
}

// returns the sales stock of the given class
func (class Class) SalesStock() Class_Stock {
	username := class.UserName
	stockList := *Users[username].ClassStocks()
	for i := 0; i < len(stockList); i++ {
		s := &stockList[i]
		if (s.Class_id == class.Id) && (s.Usage_type == `Sales`) {
			return *s
		}
	}
	return NotFoundClassStock
}

// returns the consumption stock of the given class
// under development - at present assumes there is only one
// WAS 	query := `SELECT stocks.* FROM stocks INNER JOIN commodities ON stocks.commodity_id = commodities.id where stocks.owner_id = ? AND Usage_type ="Consumption" AND commodities.name="Consumption"`
func (class Class) ConsumerGood() Class_Stock {
	username := class.UserName
	stockList := *Users[username].ClassStocks()

	for i := 0; i < len(stockList); i++ {
		s := &stockList[i]
		if (s.Class_id == class.Id) && (s.Usage_type == `Consumption`) {
			return *s
		}
	}
	return NotFoundClassStock
}

// METHODS OF INDUSTRY STOCKS

// fetches the name of the owner of this stock
func (s Industry_Stock) OwnerName() string {
	username := s.UserName
	industryList := *Users[username].Industries()
	for i := 0; i < len(industryList); i++ {
		ind := &industryList[i]
		if s.Industry_id == ind.Id {
			return ind.Name
		}
	}
	return `UNKNOWN OWNER`
}

// return the name of the commodity that the given Industry_Stock consists of
// WAS 	rows, err := db.SDB.Queryx("SELECT * FROM commodities where Id = ?", i.Commodity_id)
func (s Industry_Stock) CommodityName() string {
	username := s.UserName
	commodityList := *Users[username].Commodities()
	for i := 0; i < len(commodityList); i++ {
		c := commodityList[i]
		if s.Commodity_id == c.Id {
			return c.Name
		}
	}
	return `UNKNOWN COMMODITY`
}

// return the commodity object that the given stock consists of
// WAS 	rows, err := db.SDB.Queryx("SELECT * FROM commodities where Id = ?", i.Commodity_id)
func (s Industry_Stock) Commodity() *Commodity {
	username := s.UserName
	commodityList := *Users[username].Commodities()
	for i := 0; i < len(commodityList); i++ {
		c := commodityList[i]
		if s.Commodity_id == c.Id {
			return &c
		}
	}
	return &NotFoundCommodity
}

// under development
// will eventually be parameterised to yield value, price or quantity depending on a 'display' parameter
func (stock Industry_Stock) DisplaySize(mode string) float32 {
	switch mode {
	case `prices`:
		return stock.Size
	case `quantities`:
		return stock.Size // switch in price once this is in the model
	default:
		panic(`unknown display mode requested`)
	}
}

// (Experimental) Creates a url to link to this simulation, to be used in templates such as dashboard
// In this way all the URL naming is done in native Golang, not in the template
// We may also use such methods in the Trace function to improve usability
func (s Simulation) Link() string {
	return `/user/create/` + strconv.Itoa(s.Id)
}

// fetches the industry that owns this industry stock
// If it has none (an error, but we need to diagnose it) return nil.
func (s Industry_Stock) Industry() *Industry {
	industryList := *Users[s.UserName].Industries()
	for i := 0; i < len(industryList); i++ {
		ind := &industryList[i]
		if s.Industry_id == ind.Id {
			return ind
		}
	}
	return nil
}

// fetches the name of the industry that owns this industry stock.
// If it has none (an error, but we need to diagnose it) return "UNKNOWN INDUSTRY"
func (s Industry_Stock) IndustryName() string {
	i := s.Industry()
	if i == nil {
		return "UNKNOWN INDUSTRY"
	}
	return i.Name
}

// METHODS OF CLASS STOCKS

// fetches the class that owns this Class_stock
// If it has none (an error, but we need to diagnose it) return nil.
func (s Class_Stock) Class() *Class {
	classList := *Users[s.UserName].Classes()
	for i := 0; i < len(classList); i++ {
		ind := &classList[i]
		if s.Class_id == ind.Id {
			return ind
		}
	}
	return nil
}

// fetches the name of the Class that owns this Class_stock.
// If it has none (an error, but we need to diagnose it) return "UNKNOWN CLASS"
func (s Class_Stock) ClassName() string {
	c := s.Class()
	if c == nil {
		return "UNKNOWN CLASS"
	}
	return c.Name
}

// Return the name of the commodity that this Class_Stock consists of.
// Return "UNKNOWN COMMODITY" if this is not found.
func (s Class_Stock) CommodityName() string {
	username := s.UserName
	commodityList := *Users[username].Commodities()
	for i := 0; i < len(commodityList); i++ {
		c := commodityList[i]
		if s.Commodity_id == c.Id {
			return c.Name
		}
	}
	return `UNKNOWN COMMODITY`
}

func (u User) Get_current_state() string {
	id := u.CurrentSimulationID
	sims := *u.Simulations()
	if sims == nil {
		return "UNKNOWN"
	}

	for i := 0; i < len(sims); i++ {
		s := sims[i]
		if s.Id == id {
			return s.State
		}
	}
	return "UNKNOWN"
}

// helper function to set the state of the current simulation
// if we fail it's a programme error so we don't test for that
func (u User) Set_current_state(new_state string) {
	id := u.CurrentSimulationID
	sims := *u.Simulations()
	log.Output(1, fmt.Sprintf("resetting state to %s for user %s", new_state, u.UserName))
	for i := 0; i < len(sims); i++ {
		s := &sims[i]
		if (*s).Id == id {
			(*s).State = new_state
			return
		}
		log.Output(1, fmt.Sprintf("simulation with id %d not found", id))
	}
}

// Create a CommodityView object for display in a template
// taking data from two Commodity objects; one being viewed now,
// the other showing the state of the simulation at some time in the 'past'
func NewCommodityView(v *Commodity, c *Commodity) *CommodityView {
	newCommodityView := CommodityView{
		Id:                          v.Id,
		Name:                        v.Name,
		Origin:                      v.Origin,
		Usage:                       v.Usage,
		Size:                        Pair{Viewed: v.Size, Compared: c.Size},
		Total_Value:                 Pair{Viewed: (v.Total_Value), Compared: (c.Total_Value)},
		Total_Price:                 Pair{Viewed: (v.Total_Price), Compared: (c.Total_Price)},
		Unit_Value:                  Pair{Viewed: (v.Unit_Value), Compared: (c.Unit_Value)},
		Unit_Price:                  Pair{Viewed: (v.Unit_Price), Compared: (c.Unit_Price)},
		Turnover_Time:               Pair{Viewed: v.Turnover_Time, Compared: c.Turnover_Time},
		Demand:                      Pair{Viewed: v.Demand, Compared: c.Demand},
		Supply:                      Pair{Viewed: v.Supply, Compared: c.Supply},
		Allocation_Ratio:            Pair{Viewed: v.Allocation_Ratio, Compared: c.Allocation_Ratio},
		Monetarily_Effective_Demand: v.Monetarily_Effective_Demand,
		Investment_Proportion:       v.Investment_Proportion,
	}
	return &newCommodityView
}

func NewCommodityViews(v *[]Commodity, c *[]Commodity) *[]CommodityView {
	var newCommodityViews = make([]CommodityView, len(*v))
	for i := range *v {
		newCommodityView := NewCommodityView(&(*v)[i], &(*c)[i])
		newCommodityViews[i] = *newCommodityView
	}
	return &newCommodityViews
}
