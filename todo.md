# Todo

## Utilities

- [ ] use power readings to calculate units consumed
  - [ ] pre populate reading from last invoice
  - [x] automate the calculation
  - tenants like to argue over it, so keep both values and calc the difference
- [ ] late fees field
  - tack late fees onto bill as dollar value
- [ ] line rental charged with utilities, constant per lease
  - [ ] tack onto bill as dollar value
- [x] unit cost is global variable
- [x] due date net 14 for utility bill, global variabl
- [ ] utility invoice shows any previous unpaid invoices
- [ ] service reference number (unique per lease?)

## Rent

- [ ] Residential / Commercial rent services
- [ ] GST global variable (percentage) (commercial rent service only)
- [ ] rent cycle is per lease weekly (+6 days), fortnightly (2 x weekly), or monthly (4 x weekly)
  - [x] default to weekly
- [ ] rent paid date field (default to today)
- [ ] rent amount defaulted to lease rent variable
- [ ] rent can be paid out-of-order
  - bring list of due rent and click to pay out of order
- [ ] service reference number (unique per lease?)

## Tenant

- [ ] number plate for nx witness (could automate?)
- [ ] arbitrary notes
- [ ] multiple (arbitrary number of) contact fields

## Lease

- [ ] lease: rent bond (signup), gate key bond (signup) (static)
- [ ] select services per lease

## Misc

- [x] global: bank details
- [ ] how to handle hosted data and users
- [ ] show site first and order leases by site alphabetically
