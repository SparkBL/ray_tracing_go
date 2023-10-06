#ifndef INTERVAL
#define INTERVAL
#include <math.h>

struct interval
{
    interval(double min, double max)
    {
        this->min = min;
        this->max = max;
    }
    bool contains(double x)
    {
        return min <= x <= max;
    }
    bool surrounds(double x)
    {
        return min < x < max;
    }

private:
    double min = 0, max = 1;
};

static interval empty = interval(+INFINITY, -INFINITY);
static interval universe = interval(-INFINITY, +INFINITY);

#endif