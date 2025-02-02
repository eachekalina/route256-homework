package commands

import (
	"fmt"
	"homework/internal/app/core"
	"homework/internal/app/order"
	"os"
	"strconv"
	"text/tabwriter"
)

type OrderConsoleCommands struct {
	svc  *core.OrderCoreService
	help Command
}

func NewOrderConsoleCommands(svc *core.OrderCoreService, help Command) *OrderConsoleCommands {
	return &OrderConsoleCommands{svc: svc, help: help}
}

func (c *OrderConsoleCommands) AcceptOrderCommand(args []string) error {
	var req core.AcceptOrderRequest

	fs := createFlagSet(c.help)
	fs.Uint64Var(&req.OrderId, "order-id", 0, "specify order id")
	fs.Uint64Var(&req.CustomerId, "customer-id", 0, "specify customer id")
	fs.StringVar(&req.KeepDateString, "keep-date", "", "specify keep date")
	fs.Int64Var(&req.PriceRub, "price", 0, "specify price in rubles")
	fs.Float64Var(&req.WeightKg, "weight", 0.0, "specify weight in kg")
	fs.StringVar(&req.PackagingType, "packaging", "", "specify packaging")
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	return c.svc.AcceptOrder(req)
}

func (c *OrderConsoleCommands) ReturnOrderCommand(args []string) error {
	var id uint64

	fs := createFlagSet(c.help)
	fs.Uint64Var(&id, "order-id", 0, "specify order id")
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	return c.svc.ReturnOrder(id)
}

func (c *OrderConsoleCommands) GiveOrdersCommand(args []string) error {
	fs := createFlagSet(c.help)
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	args = fs.Args()

	ids := make([]uint64, len(args))
	for i, arg := range args {
		var err error
		ids[i], err = strconv.ParseUint(arg, 10, 64)
		if err != nil {
			return err
		}
	}
	return c.svc.GiveOrders(ids)
}

func (c *OrderConsoleCommands) ListOrdersCommand(args []string) error {
	var req core.ListOrdersRequest

	fs := createFlagSet(c.help)
	fs.Uint64Var(&req.CustomerId, "customer-id", 0, "specify customer id")
	fs.IntVar(&req.DisplayCount, "n", 0, "specify number of entries")
	fs.BoolVar(&req.FilterGiven, "stored-only", false, "display only stored orders")
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	orders, err := c.svc.ListOrders(req)
	if err != nil {
		return err
	}
	if len(orders) == 0 {
		fmt.Println("No orders are to display")
		return nil
	}
	c.printOrders(orders)
	return nil
}

func (c *OrderConsoleCommands) AcceptReturnCommand(args []string) error {
	var req core.AcceptReturnRequest

	fs := createFlagSet(c.help)
	fs.Uint64Var(&req.OrderId, "order-id", 0, "specify order id")
	fs.Uint64Var(&req.CustomerId, "customer-id", 0, "specify customer id")
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	return c.svc.AcceptReturn(req)
}

func (c *OrderConsoleCommands) ListReturnsCommand(args []string) error {
	var req core.ListReturnsRequest

	fs := createFlagSet(c.help)
	fs.IntVar(&req.Count, "n", 0, "specify number of entries")
	fs.IntVar(&req.PageNum, "page", 0, "specify page")
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	orders, err := c.svc.ListReturns(req)
	if err != nil {
		return err
	}
	if len(orders) == 0 {
		fmt.Println("No returns are available")
		return nil
	}
	c.printOrders(orders)
	return nil
}

func (c *OrderConsoleCommands) printOrders(orders []order.Order) {
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintf(
		w,
		"%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		"Order id",
		"Customer id",
		"Price RUB",
		"Weight kg",
		"Add date",
		"Keep date",
		"Is given",
		"Give date",
		"Is returned",
		"Return date")
	for _, order := range orders {
		fmt.Fprint(w, order)
	}
	w.Flush()
}
