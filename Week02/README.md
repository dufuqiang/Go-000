**Week02 作业题目：**

1. 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？



答：需要，操作数据库时，dao层遇到一个 sql.ErrNoRows 应该把错误 Wrap，抛给上层去处理。