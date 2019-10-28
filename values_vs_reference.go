package main

import (
	"fmt"
)

func main() {
	one := Inner{1}
	two := Inner{2}
	inner2_ptr := &[2]Inner{one, two}
	holder := Holder{[2]Inner{one, two}, inner2_ptr}
	print_stuff(holder)
	doesnt_change(holder)
	print_stuff(holder)
	does_change_to_100(&holder)
	print_stuff(holder)
	doesnt_change_to_200(&holder)
	print_stuff(holder)

	// ==============================
	// Now mess with the pointer list
	fmt.Println("==============================")
	print_ptr_stuff(holder)
	does_change_ptr(holder)
	print_ptr_stuff(holder)
	does_change_ptr2(&holder)
	print_ptr_stuff(holder)
	does_change_ptr3(holder)
	print_ptr_stuff(holder)
}

type Holder struct {
	mylist    [2]Inner
	mylistptr *[2]Inner
}

type Inner struct {
	myval int
}

func print_stuff(holder Holder) {
	for i, inner := range holder.mylist {
		fmt.Printf("%v, %v\n", i, inner.myval)
	}
}

func doesnt_change(holder Holder) {
	// holder passed into the function is a copy
	for i, inner := range holder.mylist {
		holder.mylist[i].myval = 100
		inner.myval = 200
	}
}

func does_change_to_100(holder *Holder) {
	// Now we have a pointer to holder, and we directly access its list via index
	for i := range holder.mylist {
		holder.mylist[i].myval = 100
	}
}

func doesnt_change_to_200(holder *Holder) {
	// We have a pointer to holder, but range returns a COPY of the value as the second argument
	for _, inner := range holder.mylist {
		// inner is still a copy! This doesn't do anything to our original holder
		inner.myval = 200
		(&inner).myval = 300
	}
}

//================================================================================
// Now we're going to change the pointer list
// It now doesn't matter if pass in holder via reference or value, because the
// copied one will still have a reference to the same list.
// However, we still CANNOT change the inner value via range

func print_ptr_stuff(holder Holder) {
	for i, inner := range holder.mylistptr {
		fmt.Printf("%v, %v\n", i, inner.myval)
	}
}

func does_change_ptr(holder Holder) {
	// holder passed into the function is a copy
	for i, inner := range holder.mylistptr {
		holder.mylistptr[i].myval = 100
		inner.myval = 200
	}
}

func does_change_ptr2(holder *Holder) {
	// holder passed into the function is a copy
	for i, inner := range holder.mylistptr {
		holder.mylistptr[i].myval = 300
		(&inner).myval = 400
	}
}

func does_change_ptr3(holder Holder) {
	for i, _ := range holder.mylistptr {
		(&holder.mylistptr[i]).myval = 500
	}
}
