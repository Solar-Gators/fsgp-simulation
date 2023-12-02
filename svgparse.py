import re
from svg.path import parse_path, Line, CubicBezier, QuadraticBezier, Arc
from scipy.interpolate import splev, splprep
import numpy as np

svg_path_data = "M518.862,312.612C518.526,317.135 516.46,332.575 514.376,336.774C512.321,340.916 504.473,345.797 500.795,347.972C418.602,396.575 336.377,445.127 254.002,493.421C225.836,509.935 171.819,540.599 142.787,555.559C119.214,567.707 94.511,579.55 67.964,583.497C60.133,584.661 53.047,584.76 49.338,577.003C40.518,558.556 55.579,539.71 66.076,525.797C85.252,500.381 103.666,472.815 129.255,453.077C142.204,443.089 159.405,435.992 169.791,422.629C189.569,397.181 204.862,363.125 239.514,355.781C265.065,350.366 287.142,370.895 310.47,377.06C327.287,381.505 343.786,373.252 358.75,366.978C370.96,361.858 383.998,356.686 392.356,345.809C401.444,333.983 404.465,318.612 410.293,305.284C415.943,292.362 427.552,282.782 440.249,277.231C462.214,267.629 487.037,276.605 506.496,287.922C514.827,292.767 520.947,302.761 518.862,312.612Z"


def curvature(dx, dy, ddx, ddy):
    return (dx * ddy - dy * ddx) / (dx**2 + dy**2) ** (3 / 2)


def get_points_and_derivatives(path_data, num_points):
    path_length = path.length()
    points = [
        path.point(i / float(num_points - 1), error=1e-5) for i in range(num_points)
    ]

    x, y = zip(*[(point.real, point.imag) for point in points])
    tck, u = splprep([x, y], s=0)

    xi, yi = splev(np.linspace(0, 1, num_points), tck)

    dx, dy = splev(np.linspace(0, 1, num_points), tck, der=1)
    ddx, ddy = splev(np.linspace(0, 1, num_points), tck, der=2)

    return xi, yi, dx, dy, ddx, ddy


def get_segment_details(path):
    segments = path if hasattr(path, "__iter__") else [path]
    segment_lengths = []

    for segment in segments:
        segment_length = segment.length(error=1e-5)

        segment_lengths.append(segment_length)

    return segment_lengths


path = parse_path(svg_path_data)

match = re.search(r'd="([^"]+)"', svg_path_data)
if match:
    path_data = match.group(1)
else:
    path_data = svg_path_data

num_points = 1000
xi, yi, dx, dy, ddx, ddy = get_points_and_derivatives(path_data, num_points)

radii_of_curvature = [
    1 / curvature(dx[i], dy[i], ddx[i], ddy[i])
    if curvature(dx[i], dy[i], ddx[i], ddy[i]) != 0
    else float("inf")
    for i in range(num_points)
]

print(radii_of_curvature)

segment_lengths = get_segment_details(path)

print(segment_lengths)
