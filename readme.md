# Avisha Fn

> Decompose a system based on volatility and implement use cases as an integeration between the areas of volatility.

## Storage

The first major source of volatility I ran into was application storage.
Storage can take many forms within a single application: in-memory, file on disk
or a remote database. Not only is the locality of data variable, but so is the API used to access it: SQL, Object Storage, Document Databases, CRUD methods, etc.
Applications often couple directly to specific storage solutions, eg "the mysql database". The entire application then orients around this decision.

On the flip side, I've seen many a project follow the "repository pattern" which
stems from the Domain Driven Design (DDD) philosophy. While repositories adhere to the idea that you should not depend on a specific storage solution, the naive implementation becomes untennable. When you start trying to encapsulate data access details you run into some interesting design questions around what your repository API should look like.

Naive functional decomposition will leave you with static query methods like `GetThingById` and `ListThingsWithThisAndThatButNotThis`. This is compounded when you have a repository interface for each important entity in the system. This makes actually implementing your interfaces tedious.

The tedium emerges because the repository is not actually encapsulating the volatility properly, which is _why_ it gets so bloated. _Tedium is a code-smell_. While it _is_ encapsulting the concrete data access, it ignores query volatiltity: the many ways you might want to query data. This is a key source of volatility for storage abstractions.

The way to solve this is with a _query api_. You want an abstraction that can capture all possible ways to compose a query. This is why some developers choose to just expose raw SQL - because SQL is literally a query language designed to solve this problem. One issue with SQL, while it does solve the problem of queries, it also couples you to using SQL compatible databases, which you may not want.

Therefore, you absolutely don't want to wrap a query language _behind_ your repository interfaces. That defeats the purpose of the repository! Use a query api or create your own as part of the repository interface.

Once you have this, you have implemented the "Read" in Create Read Update Delete (CRUD) via your query api. Now you can simply code your application use-cases to the query api, changing the queries as requirements change.

This lets you keep your repository interface small, and the application can choose how it wants to query your data.

Finally, if your language has good facilities for querying such as LINQ in C#, you don't need to implement a query api because you already have one!
