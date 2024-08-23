namespace Danom;

public static class UnitActionExtensions
{
    public static Func<TResult, Unit> ToUnitFunc<TResult>(this Action<TResult> action)
    {
        return result =>
        {
            action(result);
            return Unit.Value;
        };
    }

    public static Func<Unit, Unit> ToUnitFunc(this Action action)
    {
        return _ =>
        {
            action();
            return Unit.Value;
        };
    }
}

public static class UnitTaskExtensions
{
    public static async Task<Unit> ToUnitAsync(this Task task)
    {
        await task;
        return Unit.Value;
    }
}
