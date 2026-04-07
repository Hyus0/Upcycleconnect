import { createRouter, createWebHistory } from "vue-router";
import AdminLayout from "../components/layout/AdminLayout.vue";
import LoginPage from "../pages/LoginPage.vue";
import DashboardPage from "../pages/DashboardPage.vue";
import UsersPage from "../pages/UsersPage.vue";
import PrestationsPage from "../pages/PrestationsPage.vue";
import CategoriesPage from "../pages/CategoriesPage.vue";
import EventsPage from "../pages/EventsPage.vue";
import ModerationPage from "../pages/ModerationPage.vue";
import FinancePage from "../pages/FinancePage.vue";
import NotificationsPage from "../pages/NotificationsPage.vue";

const routes = [
  {
    path: "/login",
    name: "login",
    component: LoginPage,
    meta: { title: "Connexion", subtitle: "Acces administrateur" }
  },
  {
    path: "/",
    component: AdminLayout,
    children: [
      {
        path: "",
        name: "dashboard",
        component: DashboardPage,
        meta: { title: "Dashboard", subtitle: "Vue d'ensemble operationnelle" }
      },
      {
        path: "users",
        name: "users",
        component: UsersPage,
        meta: { title: "Utilisateurs", subtitle: "Comptes, roles et activation" }
      },
      {
        path: "prestations",
        name: "prestations",
        component: PrestationsPage,
        meta: { title: "Prestations", subtitle: "Catalogue, prix et publication" }
      },
      {
        path: "categories",
        name: "categories",
        component: CategoriesPage,
        meta: { title: "Categories", subtitle: "Taxonomie du catalogue" }
      },
      {
        path: "events",
        name: "events",
        component: EventsPage,
        meta: { title: "Evenements", subtitle: "Sessions, capacites et calendrier" }
      },
      {
        path: "moderation",
        name: "moderation",
        component: ModerationPage,
        meta: { title: "Moderation", subtitle: "Validation et archivage des contenus" }
      },
      {
        path: "finance",
        name: "finance",
        component: FinancePage,
        meta: { title: "Finances", subtitle: "Suivi des encaissements et engagements" }
      },
      {
        path: "notifications",
        name: "notifications",
        component: NotificationsPage,
        meta: { title: "Notifications", subtitle: "Campagnes et messages transactionnels" }
      }
    ]
  }
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
});

router.afterEach((to) => {
  document.title = `UpcycleConnect | ${to.meta.title ?? "Admin"}`;
});

export default router;
