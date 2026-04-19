import { createRouter, createWebHistory } from "vue-router";

import Register from "./pages/Register.vue";
import Login from "./pages/Login.vue";
import Profil from "./pages/Profil.vue";

import Home from "./pages/Home.vue";
import Annonces from "./pages/Annonces.vue";

const routes = [
  { path: "/", component: Home },
  { path: "/inscription", component: Register },
  { path: "/connexion", component: Login },
  { path: "/annonces", component: Annonces },
  {
    path: "/profil",
    component: Profil,
    meta: { requiresAuth: true },
    children: [
      {
        path: "",
        name: "dashboard",
        component: () => import("./pages/profilpages/DashboardHome.vue"),
      },
      {
        path: "annonces",
        name: "mes-annonces",
        component: () => import("./pages/profilpages/Annonces.vue"),
      },
      {
        path: "depots",
        name: "mes-depots",
        component: () => import("./pages/profilpages/Depots.vue"),
      },
      {
        path: "informations",
        name: "mes-informations",
        component: () => import("./pages/profilpages/Infos.vue"),
      },
      {
        path: "password",
        name: "modification-password",
        component: () => import("./pages/profilpages/ChangePassword.vue"),
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (to.hash) {
      return { el: to.hash, behavior: "smooth" };
    }
    return { top: 0 };
  },
});

router.beforeEach(async (to, from, next) => {
  if (to.matched.some((record) => record.meta.requiresAuth)) {
    const id = localStorage.getItem("userId");
    const token = localStorage.getItem("userToken");

    if (!id || !token) {
      return next("/connexion");
    }

    try {
      const response = await fetch(
        `http://localhost:8081/check-session?id=${id}`,
        {
          headers: { Authorization: token },
        },
      );
      const data = await response.json();

      if (data.isValid) {
        next();
      } else {
        localStorage.clear();
        next("/connexion");
      }
    } catch (error) {
      next("/");
    }
  } else {
    next();
  }
});

export default router;
