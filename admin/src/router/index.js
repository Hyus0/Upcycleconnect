import { createRouter, createWebHistory } from "vue-router";
import AdminLayout from "../components/layout/AdminLayout.vue";
import DashboardPage from "../pages/DashboardPage.vue";
import UsersPage from "../pages/UsersPage.vue";
import PrestationsPage from "../pages/PrestationsPage.vue";
import CategoriesPage from "../pages/CategoriesPage.vue";
import EventsPage from "../pages/EventsPage.vue";

const routes = [
  {
    path: "/",
    component: AdminLayout,
    children: [
      { path: "", name: "dashboard", component: DashboardPage, meta: { title: "Dashboard" } },
      { path: "users", name: "users", component: UsersPage, meta: { title: "Utilisateurs" } },
      {
        path: "prestations",
        name: "prestations",
        component: PrestationsPage,
        meta: { title: "Prestations" }
      },
      {
        path: "categories",
        name: "categories",
        component: CategoriesPage,
        meta: { title: "Categories" }
      },
      { path: "events", name: "events", component: EventsPage, meta: { title: "Evenements" } }
    ]
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.afterEach((to) => {
  document.title = `UpcycleConnect | ${to.meta.title ?? "Admin"}`;
});

export default router;
