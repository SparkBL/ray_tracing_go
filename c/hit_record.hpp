#ifndef HIT_RECORD
#define HIT_RECORD
#include "vector.hpp"
#include "ray.hpp"
#include <memory>

using std::shared_ptr

    class material;

struct hit_record
{
    point p;
    vec3 normal;
    shared_ptr<material> mat;
    double t;
    bool is_front_face;

    void set_face_normal(const ray &r, const vec3 &outward_normal)
    {
        // Sets the hit record normal vector.
        // NOTE: the parameter `outward_normal` is assumed to have unit length.

        is_front_face = (r.direction && outward_normal) < 0;
        if (is_front_face)
        {
            normal = outward_normal;
        }
        else
        {
            normal = -outward_normal;
        }
    }
};

struct scatter_record
{
    color attenuation;
    ray scattered;
    bool is_scattered;
    scatter_record(bool is_scattered, ray scattered, color attenuation)
    {
        this->attenuation = attenuation;
        this->scattered = scattered;
        this->is_scattered = is_scattered;
    }
};

class material
{
public:
    virtual scatter_record scatter(const ray r_in, const hit_record &rec) const = 0;
    virtual ~material() = default;
};

class lambertian : public material
{
private:
    color albedo;

public:
    scatter_record scatter(const ray r_in, const hit_record &rec) const override
    {
        vec3 scatter_direction = rec.normal + random_unit_vector();
        if (scatter_direction.is_close_to_zero())
        {
            scatter_direction = rec.normal;
        }
        return scatter_record(true, ray(rec.p, scatter_direction), albedo);
    }
};

#endif