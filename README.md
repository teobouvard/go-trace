# gotrace

`gotrace` is a Golang implementation of Peter Shirley's [excellent books on raytracing](https://raytracing.github.io/).

I used it as a way to better understand Go interfaces. Being more familiar with object-oriented languages, I wondered how polymorphism behaviour could be implemented without inheritance. Turns out it is pretty awesome !

![render of the second book cover](assets/final_scene.jpg)
![render of the first book cover](assets/cornell_box.jpg)

## Pros

- Strong concurrency primitives. Using weighted semaphores is a very intuitive way to create a workgroup.
- Interfaces feels like more natural way to provide polymorphism behaviour at runtime, compared to inheritance.
- Awesome [self-documentation](https://pkg.go.dev/github.com/teobouvard/gotrace)

## Cons

- No operator overloading

```go
corner := lookFrom.Sub(u.Scale(width * focusDist)).Sub(v.Scale(height * focusDist)).Sub(w.Scale(focusDist))
```

would have been clearer as

```go
corner := lookfrom - u*width*focusDist - v*height*focusDist - w*focusDist
```

- No forward declarations

`Actor` needs to know about `Shape`, `Material` and `Ray`, but `Shape` also has to know about `Ray`, so they can't live in different packages without having to add unecessary complexity, because Go doesn't support forward declarations and can't resolve "circular" dependencies.

## Credits

- [NASA Blue Marble](https://visibleearth.nasa.gov/collection/1484/blue-marble)
- [Peter Shirley's original work](https://raytracing.github.io/)
