#ifndef HIT_RECORD
#define HIT_RECORD
#include "vector.hpp"
#include "ray.hpp"
#include <memory>

using std::shared_ptr;

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

    scatter_record(bool is_scattered, const ray &scattered, const color &attenuation)
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

class metal : public material
{
private:
    color albedo;
    double fuzziness;

public:
    scatter_record scatter(const ray r_in, const hit_record &rec) const override
    {
        vec3 reflected = reflect(unit_vector(r_in.direction), rec.normal);
        ray scattered = ray(rec.p, reflected + random_unit_vector() * fuzziness);
        return scatter_record(true, scattered, albedo);
    }
};

class dielectric : public material
{
private:
    double refraction_index;

    double reflectance(double cosine, double ref_idx) const
    {
        double r0 = (1 - ref_idx) / (1 + ref_idx);
        r0 = r0 * r0;
        return r0 + (1 - r0) * pow(1 - cosine, 5);
    }

public:
    scatter_record scatter(const ray r_in, const hit_record &rec) const override
    {
        color attenuation = color(1, 1, 1);
        double refraction_ratio = refraction_index;
        if (rec.is_front_face)
        {
            refraction_ratio = 1.0 / refraction_index;
        }
        vec3 unit_direction = unit_vector(r_in.direction);
        double cos_theta = fmin(-unit_direction && rec.normal, 1.0);
        double sin_theta = sqrt(1.0 - cos_theta * cos_theta);
        vec3 direction;
        bool cannot_refract = refraction_ratio * sin_theta > 1.0;
        if (cannot_refract || reflectance(cos_theta, refraction_ratio) > dis(gen))
        {
            direction = reflect(unit_direction, rec.normal);
        }
        else
        {
            direction = refract(unit_direction, rec.normal, refraction_ratio);
        }
        ray scattered = ray(rec.p, direction);
        return scatter_record(true, scattered, attenuation);
    }
};
#endif