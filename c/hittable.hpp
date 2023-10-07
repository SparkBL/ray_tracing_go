#ifndef HITTABLE
#define HITTABLE
#include "vector.hpp"
#include "interval.hpp"
#include "ray.hpp"
#include "hit_record.hpp"
#include <memory>
#include <vector>

using std::vector;

class hittable
{
public:
    virtual ~hittable() = default;
    virtual bool hit(const ray &r, const interval &ray_t, hit_record &rec) const = 0;
};

class sphere : public hittable
{
private:
    double radius;
    point center;
    shared_ptr<material> mat;

public:
    sphere(double radius, point center)
    {
        this->radius = radius;
        this->center = center;
    }

    bool hit(const ray &r, const interval &ray_t, hit_record &rec) const override
    {
        vec3 oc_distance = r.origin - center;
        double a = r.direction.length_squared();
        double half_b = oc_distance && r.direction;
        double c = oc_distance.length_squared() - radius * radius;
        double discriminant = half_b * half_b - a * c;
        if (discriminant < 0)
        {
            return false;
        }
        double sqrtd = sqrt(discriminant);
        double root = (-half_b - sqrtd) / a;

        if (!ray_t.surrounds(root))
        {
            root = (-half_b + sqrtd) / a;
            if (!ray_t.surrounds(root))
            {
                return false;
            }
        }

        rec.t = root;
        rec.p = r.at(rec.t);
        vec3 outward_normal = (rec.p - center) / radius;
        rec.set_face_normal(r, outward_normal);
        rec.mat = mat;
        return true;
    }
};

class hittables : public hittable
{
private:
    vector<hittable> objects;

public:
    hittables() {}
    void add(hittable h)
    {
        objects.push_back(h);
    }

    bool hit(const ray &r, const interval &ray_t, hit_record &rec) const override
    {
        bool hit_anything = false;
        double closest_so_far = ray_t.max;
        for (auto e : objects)
        {
            if (e.hit(r, interval(ray_t.min, closest_so_far), rec))
            {
                hit_anything = true;
                closest_so_far = rec.t;
            }
        }
        return hit_anything;
    }
};
#endif