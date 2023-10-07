#ifndef CAMERA
#define CAMERA
#include "vector.hpp"
#include "ray.hpp"
#include "hittable.hpp"
#include "hit_record.hpp"

struct camera
{
private:
    double aspect_ratio;
    int image_width, image_height;
    point center;
    point pixel_zero_location;
    vec3 pixel_delta_u, pixel_delta_v;

    int samples_per_pixel;
    int max_ray_depth;

    double focal_length;
    double viewport_height, viewport_width;

    vec3 pixel_sample_square()
    {
        double px = -0.5 * dis(gen), py = -0.5 * dis(gen);
        return pixel_delta_u * px + pixel_delta_v * py;
    }

    ray get_ray(int i, int j)
    {
        point pixel_center = pixel_zero_location +
                             (pixel_delta_u * i) +
                             (pixel_delta_v * j);
        point pixel_sample = pixel_center + pixel_sample_square();
        return ray(center, pixel_sample - center);
    }

    color ray_color(const ray *r, int depth, hittable world)
    {
        hit_record rec;
        if (depth < 0)
        {
            return color(0, 0, 0);
        }
    }

public:
    camera()
    {
        aspect_ratio = 16.0 / 9.0;
        image_width = 800;
        image_height = static_cast<int>(image_width / aspect_ratio);

        focal_length = 1.0;
        viewport_height = 2.0;
        viewport_width = viewport_height * (image_width / image_height);

        center = point(0, 0, 0);

        samples_per_pixel = 10;
        max_ray_depth = 50;

        if (image_height < 1)
        {
            image_height = 1;
        }

        vec3 viewport_u = vec3(viewport_width, 0, 0);
        vec3 viewport_v = vec3(0, -viewport_height, 0);

        pixel_delta_u = viewport_u / image_width;
        pixel_delta_v = viewport_v / image_height;

        vec3 viewport_up_left = center +
                                vec3(0, 0, -focal_length) -
                                (viewport_u / 2) -
                                (viewport_v / 2);
        pixel_zero_location = viewport_up_left + (pixel_delta_u + pixel_delta_v) * 5;
    }
};

#endif