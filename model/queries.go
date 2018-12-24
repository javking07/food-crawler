package model

//FoodDataCreateTableQuery is sql query for creating fda_data table
const SqlFoodDataCreateTableQuery string = `CREATE TABLE IF NOT EXISTS fda_data
(
id SERIAL,
data JSON NOT NULL,
name TEXT NOT NULL,
kind TEXT NOT NULL,
color TEXT NOT NULL,
age INT NOT NULL,
CONSTRAINT fda_data_pkey PRIMARY KEY (id)
)`

//TestFoodDataCreateTableQuery is sql query for creating fda_data table
const SqlTestFoodDataCreateTableQuery string = `use testspace; create table govdata 
(
id UUID PRIMARY KEY,
data blob,
);`
