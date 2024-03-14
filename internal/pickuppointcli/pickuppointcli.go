package pickuppointcli

import (
	"bufio"
	"errors"
	"fmt"
	"homework/internal/model"
	service "homework/internal/pickuppointservice"
	"homework/internal/storage"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

// PickUpPointCli provides a console line interface for working with pick-up points storage.
type PickUpPointCli struct {
	serv    *service.PickUpPointService
	scanner *bufio.Scanner
}

// NewPickUpPointCli creates a new PickUpPointCli
func NewPickUpPointCli(storagepath string) (*PickUpPointCli, error) {
	stor, err := storage.NewPickUpPointFileStorage(storagepath)
	if err != nil {
		return nil, err
	}
	serv := service.NewPickUpPointService(&stor)
	return &PickUpPointCli{serv: serv, scanner: bufio.NewScanner(os.Stdin)}, nil
}

// Run starts the console line interface
func (c *PickUpPointCli) Run() error {
	defer c.serv.Close()
	for {
		exit, err := c.handleCommand()
		if err != nil {
			log.Println(err)
		}
		if exit {
			return nil
		}
	}
}

func (c *PickUpPointCli) getLine(prompt string) (string, bool, error) {
	fmt.Printf("%s > ", prompt)
	if !c.scanner.Scan() {
		return "", true, c.scanner.Err()
	}
	line := c.scanner.Text()
	return strings.TrimSpace(line), false, nil
}

func (c *PickUpPointCli) getUint(prompt string) (uint64, bool, error) {
	str, exit, err := c.getLine(prompt)
	if err != nil || exit {
		return 0, exit, err
	}
	n, err := strconv.ParseUint(str, 10, 64)
	return n, false, err
}

func (c *PickUpPointCli) handleCommand() (bool, error) {
	line, exit, err := c.getLine("")
	if err != nil || exit {
		return exit, err
	}
	switch line {
	case "help":
		c.printHelp()
		return false, nil
	case "exit":
		return true, nil
	case "create":
		return c.handleCreate()
	case "list":
		c.handleList()
		return false, nil
	case "get":
		return c.handleGet()
	case "update":
		return c.handleUpdate()
	case "delete":
		return c.handleDelete()
	default:
		return false, errors.New("no such command found, use `help` for help")
	}
}

func (c *PickUpPointCli) printHelp() {
	fmt.Println(`Available commands:

help	show this help
exit	exit program
create	create a new pick-up point
list	list all existing pick-up points
get		show selected pick-up point
update	update pick-up point
delete	delete pick-up point`)
}

func (c *PickUpPointCli) handleCreate() (bool, error) {
	id, exit, err := c.getUint("enter point id")
	if err != nil || exit {
		return exit, err
	}
	name, exit, err := c.getLine("enter point name")
	if err != nil || exit {
		return exit, err
	}
	address, exit, err := c.getLine("enter point address")
	if err != nil || exit {
		return exit, err
	}
	contact, exit, err := c.getLine("enter point contact")
	if err != nil || exit {
		return exit, err
	}
	point := model.PickUpPoint{
		Id:      id,
		Name:    name,
		Address: address,
		Contact: contact,
	}
	return false, c.serv.CreatePoint(point)
}

func (c *PickUpPointCli) handleList() {
	points := c.serv.ListPoints()
	printPoints(points)
}

func (c *PickUpPointCli) handleGet() (bool, error) {
	id, exit, err := c.getUint("enter point id")
	if err != nil || exit {
		return exit, err
	}
	point, err := c.serv.GetPoint(id)
	if err != nil {
		return false, err
	}
	printPoints([]model.PickUpPoint{point})
	return false, nil
}

func (c *PickUpPointCli) handleUpdate() (bool, error) {
	id, exit, err := c.getUint("enter point id")
	if err != nil || exit {
		return exit, err
	}
	name, exit, err := c.getLine("enter new point name")
	if err != nil || exit {
		return exit, err
	}
	address, exit, err := c.getLine("enter new point address")
	if err != nil || exit {
		return exit, err
	}
	contact, exit, err := c.getLine("enter new point contact")
	if err != nil || exit {
		return exit, err
	}
	point := model.PickUpPoint{
		Id:      id,
		Name:    name,
		Address: address,
		Contact: contact,
	}
	return false, c.serv.UpdatePoint(point)
}

func (c *PickUpPointCli) handleDelete() (bool, error) {
	id, exit, err := c.getUint("enter point id")
	if err != nil || exit {
		return exit, err
	}
	return false, c.serv.DeletePoint(id)
}

func printPoints(points []model.PickUpPoint) {
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintf(
		w,
		"%s\t%s\t%s\t%s\n",
		"Id",
		"Name",
		"Address",
		"Contact")
	for _, point := range points {
		fmt.Fprint(w, point)
	}
	w.Flush()
}
