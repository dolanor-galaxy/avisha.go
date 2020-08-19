<!-- markdownlint-disable MD029 -->

# Avisha

> Decompose a system based on volatility and implement use cases as an integeration between the areas of volatility.

## Storage

### Volatility

The first major source of volatility I ran into was application storage.
Storage can take many forms within a single application: in-memory, file on disk
or a remote database. Not only is the locality of data variable, but so is the API used to access it: SQL, Object Storage, Document Databases, CRUD methods, etc.
Applications often couple directly to specific storage solutions, eg "_the_ mysql database". The entire application then orients around this decision.

On the flip side, I've seen many a project follow the "repository pattern" which
stems from the Domain Driven Design (DDD) philosophy. While repositories adhere to the idea that you should not depend on a specific storage solution, the naive implementation becomes untennable. When you start trying to encapsulate data access details you run into some interesting design questions around what your repository API should look like.

Naive functional decomposition will leave you with static query methods like `GetThingById` and `ListThingsWithThisAndThatButNotThis`. This is compounded when you have a repository interface for each important entity in the system. This makes implementing the interfaces tedious.

The tedium emerges because the repository is not actually encapsulating the volatility properly, which is _why_ it gets so bloated. _Tedium is a code-smell_. While it _is_ encapsulting the concrete data access, it ignores query volatiltity: the many ways you might want to query data. This is a key source of volatility for storage abstractions.

The way to solve this is with a _query api_. You want an abstraction that can capture all possible ways to compose a query. This is why some developers choose to just expose raw SQL - because SQL is literally a query language designed to solve this problem. One issue with SQL, while it does solve the problem of queries, it also couples you to using SQL compatible databases, which you may not want.

> Do not wrap a query language _behind_ your repository interfaces. This defeats the purpose of the repository! Use a query api or create your own as part of the repository interface.

Once you have this, you have implemented the "Read" in Create Read Update Delete (CRUD) via your query api. Now you can simply code your application use-cases to the query api, changing the queries as requirements change.

This lets you keep your repository interface small, and the application can choose how it wants to query your data.

Finally, if your language has good facilities for querying such as LINQ in C#, you don't need to implement a query api because you already have one!

### Generics and the Dynamism of Data Storage

Data is typically encoded in a format suitable for the storage context. It doesn't just stay in the encoding of the programming language being used: you can't just use the data from disk "as is", and start using it at the application level. For example SQL databases don't have "structs" or "objects", they have "tables". Neither the data layout, encoding or logical relationships directly map onto the application entities. There is some overlap of course, hence ORMs (object-relational mapper). It is not uncommon to build your entities using data from multiple tables. This is one area where the ORM abstraction breaks down. Storing data object-wise, while convenient, removes the flexibility and precision at the heart of SQL and normalized data.

In princple you will always need logic that maps arbitrary data from a storage system to application entities. This logic is your "query" logic. On the application level the query api needs to be composable so as to capture any possible permutation of a query it might want. The storage driver consumes the query api and maps it to the storage system in question. Querying storage systems is inherently volatile as queries change over time.

Initially, when exploring how to write a query api in Go, I ran into what I thought was a show-stopper problem with the language: _no generics_. In reality, it could be a problem but it wont stop the show. Let me explain.

<!-- Repository interfaces either require a static definition for each entity (`TenantRepository`, `SiteRepository`, `LeaseRepository`) with CRUD methods and a corresponding implementation for each storage object, or a single generic interface that handles all entities and therefore only one implementation for each storage object. -->

The application wants to query the data storage and extract concrete types from it. There are a number of stragies you can deploy to satsify this:

1. Static repository interface for each concrete type, and a static method for each query.

```go
type TenantRepository interface {
    GetByName(string) (Tenant, bool)
    ListAll() []Tenant
    ListByContactMethod(ContactMethod) []Tenant
    ListIfContactIsNotEmpty() []Tenant
    // ...
}

type SiteRepository interface {
    GetByNumber(string) (Site, bool)
    ListAll() []Site
    ListByDwelling(Dwelling) []Site
    ListVacancies() []Site
    // ...
}

type LeaseRepository interface {
    Get(tenant string, site string) (Lease, bool)
    ListAll() []Lease
    ListExpired() []Lease
    ListActive() []Lease
    ListByTenant(string) []Lease
    ListBySite(string) []Lease
    // ...
}
```

2. Static interface for each concrete type, and a single query method each.

```go
type TenantRepository interface {
    ListTenant([]func(Tenant) bool) []Tenant
}

type SiteRepository interface {
    ListSite([]func(Site) bool) []Site
}

type LeaseRepository interface {
    ListLease([]func(Lease) bool) []Lease
}

type Repository interface {
    TenantRepository
    SiteRepository
    LeaseRepository
}
```

3. Dynamic interface for any type, and a single query method.

```go
type Repository interface {
    List([]func(interface{}) bool) []interface{}
}
```

4. Generic interface for any type, and a single query method.

```go
type Repository(type T) interface {
    List([]func(T) bool) []T
}
```

As you can see, each sample the code gets progressively dry. To me, it's obvious that capturing the query volatility provides the most value, thus option 1 is clearly to be avoided. Option 2 and 3 tradeoff between ease of implementation and ease of consumption. Option 2 is easier to consume, while more effort to implement since you have to duplicate code structure for each type, and option 3 is harder to consume but less tedious to implement because you only implement it once. Option 4 is preferable, but not possible in Go at the moment. It trades off compile-time by letting the compiler handle code duplication per concrete type.

Since Go doesn't have generics (yet), how can you create a generic repository?
Well, you _can_ have a flexible api using `interface{}` and thus employ runtime reflection. This comes at both a CPU cost and programmer cost. Using the reflection api places a huge burden on the programmer in terms of correctness, maintainability and morale, on top of the performance hit. It also hurts ergonomics, since the consumer of the api needs to type assert the result of every query.

This puts you in an awkward position where you're writing more code for less performance which could be considered an "anti-optimisation". Usually you write more, complex code, to result in _better_ runtime performance not reduce it! However, it is still _less_ tedious than stamping out identical repositories for each entity.

The benefit of compile-time generics in this context is less about solving the problem of flexible queries, but about making it correct, ergonomic and performant through compiler automation.

Consider the following snippets:

Go

```go
// hasContact ensures the tenant has a non-zero contact field.
hasContact := func(ent interface{}) bool {
    if tenant, ok := ent.(*Tenant); ok {
        return tenant.Contact != ""
    }
    return false
}


var tenants []*Tenant

// Query the storage and collect the results.
// Because of reflection, the result is an interface{}([]interface{}) which
// makes using it rather tedious to use.
results := storage.List(hasContact)
for _, r := range results {
    if t, ok := r.(*Tenant); ok {
        tenants = append(tenants, t)
    }
}

fmt.Printf("%#v\n", tenants)
```

Rust

```rust
// has_contact ensures the tenant has a non-zero contact field.
let has_contact = |t: &Tenant| -> bool { !t.contact.is_empty() };

// Query the storage and collect results.
let tenants = storage.list::<Tenant>(has_contact).collect::<Vec<_>>();

println!("{:#}", tenants);
```

As you can see, consuming the generic api is _much_ nicer. Let's try Go with some rustic generics.

Rustic Go

```Go
// hasContact ensures the tenant has a non-zero contact field.
hasContact := func(t Tenant) bool {
    return tenant.Contact != ""
}

// Query the storage and collect the results.
tenants := storage.List::<Tenant>(hasContact)

fmt.Printf("%#v\n", tenants)
```

To me the benefit of generics is clear. Using volatility analysis, generics captures the volatility of code structure. It lets the compiler stamp out all the type permutations of an algorithm by paramatarising the algorithm by each conrete type. In Go, the programmer has to either do this manually by hand, dynamically through reflection, or include extra tooling (code generation) to achieve this.
