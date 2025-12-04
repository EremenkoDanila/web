export const ROUTES = {
  HOME: "/",
  SOFTWARE: "/software",
  SOFTWARE_DETAIL: "/software/:id",
} as const;

export type RouteKeyType = keyof typeof ROUTES;

export const ROUTE_LABELS: { [key in RouteKeyType]: string } = {
  HOME: "Главная",
  SOFTWARE: "Программное обеспечение",
  SOFTWARE_DETAIL: "Программа",
};