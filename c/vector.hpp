#ifndef VEC3
#define VEC3
#include <math.h>
#include <sstream>
#include <random>

using std::mt19937;
using std::sqrt;
using std::string;
using std::stringstream;
using std::uniform_real_distribution;

std::random_device rd;  // Will be used to obtain a seed for the random number engine
std::mt19937 gen(rd()); // Standard mersenne_twister_engine seeded with rd()
std::uniform_real_distribution<> dis(0.0, 0.99999);

struct vec3
{
    vec3() {}

    vec3(double e0, double e1, double e2)
    {
        e[0] = e0;
        e[1] = e1;
        e[2] = e2;
    }

    double x() const
    {
        return e[0];
    }

    double y() const
    {
        return e[1];
    }

    double z() const
    {
        return e[2];
    }

    vec3 operator-() const
    {
        return vec3(-e[0], -e[1], -e[2]);
    }

    vec3 operator+(const vec3 &v) const
    {
        return vec3(e[0] + v.e[0], e[1] + v.e[1], e[2] + v.e[2]);
    }

    vec3 operator-(const vec3 &a) const
    {
        return *this + (-a);
    }

    vec3 operator*(const vec3 &v) const
    {
        return vec3(e[0] * v.e[0], e[1] * v.e[1], e[2] * v.e[2]);
    }

    vec3 operator/(const vec3 &v) const
    {
        return vec3(e[0] * v.e[0], e[1] * v.e[1], e[2] * v.e[2]);
    }

    vec3 operator+(const double a) const
    {
        return vec3(e[0] + a, e[1] + a, e[2] + a);
    }

    vec3 operator*(const double a) const
    {
        return vec3(e[0] * a, e[1] * a, e[2] * a);
    }

    vec3 operator/(const double a) const
    {
        return vec3(e[0] / a, e[1] / a, e[2] / a);
    }

    double operator[](const int &idx) const
    {
        if (idx < 3)
        {
            return e[idx];
        }
        return 0;
    }

    // dor product
    double operator&&(const vec3 &v) const
    {
        return e[0] * v[0] +
               e[1] * v[1] +
               e[2] * v[2];
    }

    // cross product
    vec3 operator||(const vec3 &v) const
    {
        return vec3(
            e[1] * v[2] - e[2] * v[1],
            e[2] * v[0] - e[0] * v[2],
            e[0] * v[1] - e[1] * v[0]);
    }

    double length_squared() const
    {
        return e[0] * e[0] + e[1] * e[1] + e[2] * e[2];
    }

    double length() const
    {
        return sqrt(length_squared());
    }

    bool is_close_to_zero() const
    {
        return fabs(e[0]) < 1e-8 && fabs(e[1]) < 1e-8 && fabs(e[2]) < 1e-8;
    }

    string string() const
    {
        stringstream ss;
        ss << e[0] << ' ' << e[1] << ' ' << e[2];
        return ss.str();
    }

private:
    double e[3] = {0, 0, 0};
};

typedef vec3 color;
typedef vec3 point;

vec3 unit_vector(const vec3 &v)
{
    return v / v.length();
}

vec3 random_vector()
{
    return vec3(dis(gen), dis(gen), dis(gen));
}

vec3 random_vector(double min, double max)
{
    return vec3(
        min + dis(gen) * (max - min),
        min + dis(gen) * (max - min),
        min + dis(gen) * (max - min));
}

vec3 random_in_unit_sphere()
{
    while (true)
    {
        vec3 v = random_vector(-1, 1);
        if (v.length_squared() < 1)
        {
            return v;
        }
    }
}

vec3 random_unit_vector()
{
    return unit_vector(random_in_unit_sphere());
}

vec3 random_on_hemishere(const vec3 &normal)
{
    vec3 ushp = random_unit_vector();
    if ((ushp && normal) > 0.0)
    {
        return ushp;
    }
    else
    {
        return -ushp;
    }
}

vec3 reflect(const vec3 &v, const vec3 &n)
{
    return v - (n * (v && n));
}

vec3 refract(const vec3 &v, const vec3 &n, double et)
{
    double cos_theta = fmin(-v && n, 1.0);
    vec3 r_out_perp = (v + n * cos_theta) * et;
    vec3 r_out_parallel = n * -sqrt(abs(1 - r_out_perp.length_squared()));
    return r_out_perp + r_out_parallel;
}

#endif