package model

import "encoding/json"

func (pay Payment) MarshalBinary() ([]byte, error) {
	return json.Marshal(pay)
}
func (pay *Payment) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &pay)
}

func (i Item) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}
func (i *Item) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &i)
}

func (di DeliveryItem) MarshalBinary() ([]byte, error) {
	return json.Marshal(di)
}
func (di *DeliveryItem) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &di)
}

func (d Delivery) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}
func (d *Delivery) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &d)
}

func (c Customer) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}
func (c *Customer) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &c)
}

func (o Order) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}
func (o *Order) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (r Reply) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}
func (r *Reply) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &r)
}
