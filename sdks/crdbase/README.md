# sealos CRDBase

A lightweight `Data Store` using kubernetes's `Custom Resource`(CRD), support Traditional CURDs like Set, Update, List, Delete, along with index builtin functional.

It is design to store and query small-to-medium size of pieces of persist storage like traditional NoSQL database but keep no more requirements, just using `kubernetes`'s crd and etcd.

Currently, it must be used alone with controller runtime, and is not a standalone database.

## Feature

crdbase can provide 

1. Auto generate `CRD` and `CR` for your model struct.
2. Auto generate `Index` based on struct tags.
3. Can set `CR` owner for manage.





## Usage


### Cleanup


## Benchmark

// TODO

## TODO List
1. 字段的Unique的实现
2. 事物是否支持
3. 多表联查/relation同步
