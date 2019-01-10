package model

//FoodDataCreateTableQuery is sql query for creating fda_data table
const SqlFoodDataCreateTableQuery string = `CREATE TABLE IF NOT EXISTS fda_data
(
id UUID PRIMARY KEY,
name TEXT NOT NULL,
data blob,
)`

//TestFoodDataCreateTableQuery is sql query for creating fda_data table
const SqlTestFoodDataCreateTableQuery string = `use testspace; create table govdata 
(
id UUID PRIMARY KEY,
name TEXT NOT NULL,
data blob,
);`
