#ifndef RAY
#define RAY
#include <math.h>
#include "vector.hpp"

struct ray
{
    ray() {}

    ray(const point &origin, const vec3 &direction)
    {
        this->origin = origin;
        this->direction = direction;
    }

    point at(double t) const
    {
        return origin + direction * t;
    }

    // private:
    point origin;
    vec3 direction;
};

#endif