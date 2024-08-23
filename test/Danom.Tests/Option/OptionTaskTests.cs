namespace Danom.Tests;

using Xunit;

public sealed class OptionAsyncTests
{
    [Fact]
    public async Task MatchShouldWork()
    {
        var optionSome = await Option.SomeAsync(1).MatchAsync(x => x, () => -1);
        Assert.Equal(1, optionSome);

        var optionNone = await Option<int>.NoneAsync().MatchAsync(_ => -1, () => 1);
        Assert.Equal(1, optionNone);
    }

    [Fact]
    public async Task BindShouldWork()
    {
        AssertOption.IsSome(2, await Option.SomeAsync(1).BindAsync(x => Option.Some(x + 1)));
        AssertOption.IsSome(2, await Option.SomeAsync(1).BindAsync(x => Option.SomeAsync(x + 1)));
        AssertOption.IsNone(await Option<int>.NoneAsync().BindAsync(x => Option.Some(x + 1)));
        AssertOption.IsNone(await Option<int>.NoneAsync().BindAsync(x => Option.SomeAsync(x + 1)));
    }

    [Fact]
    public async Task MapShouldWork()
    {
        AssertOption.IsSome(2, await Option.SomeAsync(1).MapAsync(x => x + 1));
        AssertOption.IsSome(2, await Option.SomeAsync(1).MapAsync(x => Task.FromResult(x + 1)));
        AssertOption.IsNone(await Option<int>.NoneAsync().MapAsync(x => x + 1));
        AssertOption.IsNone(await Option<int>.NoneAsync().MapAsync(x => Task.FromResult(x + 1)));
    }

    [Fact]
    public async Task DefaultValueShouldWork()
    {
        Assert.Equal(1, await Option<int>.NoneAsync().DefaultValueAsync(1));
        Assert.Equal(1, await Option<int>.NoneAsync().DefaultValueAsync(Task.FromResult(1)));
        Assert.Equal(2, await Option<int>.SomeAsync(2).DefaultValueAsync(1));
        Assert.Equal(2, await Option<int>.SomeAsync(2).DefaultValueAsync(Task.FromResult(1)));
    }

    [Fact]
    public async Task DefaultWithShouldWork()
    {
        Assert.Equal(1, await Option<int>.NoneAsync().DefaultWithAsync(() => 1));
        Assert.Equal(1, await Option<int>.NoneAsync().DefaultWithAsync(() => Task.FromResult(1)));
        Assert.Equal(2, await Option<int>.SomeAsync(2).DefaultWithAsync(() => 1));
        Assert.Equal(2, await Option<int>.SomeAsync(2).DefaultWithAsync(() => Task.FromResult(1)));
    }

    [Fact]
    public async Task OrElseShouldWork()
    {
        AssertOption.IsSome(1, await Option<int>.NoneAsync().OrElseAsync(Option.Some(1)));
        AssertOption.IsSome(1, await Option<int>.NoneAsync().OrElseAsync(Option.SomeAsync(1)));
        AssertOption.IsSome(2, await Option.SomeAsync(2).OrElseAsync(Option.Some(1)));
        AssertOption.IsSome(2, await Option.SomeAsync(2).OrElseAsync(Option.SomeAsync(1)));
    }

    [Fact]
    public async Task OrElseWithShouldWork()
    {
        AssertOption.IsSome(1, await Option<int>.NoneAsync().OrElseWithAsync(() => Option.Some(1)));
        AssertOption.IsSome(1, await Option<int>.NoneAsync().OrElseWithAsync(() => Option.SomeAsync(1)));
        AssertOption.IsSome(2, await Option.SomeAsync(2).OrElseWithAsync(() => Option.Some(1)));
        AssertOption.IsSome(2, await Option.SomeAsync(2).OrElseWithAsync(() => Option.SomeAsync(1)));
    }
}
