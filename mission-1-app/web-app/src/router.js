import { createRouter, createWebHistory } from "vue-router";

import Register from "./pages/Register.vue";
import Login from "./pages/Login.vue";
import Profil from "./pages/Profil.vue";

import Home from "./pages/Home.vue";
import Annonces from "./pages/Annonces.vue";
import Forums from "./pages/Forums.vue";

import Formations from "./pages/Formations.vue"
import Evenements from "./pages/Evenements.vue"

const routes = [
  { path: "/", component: Home },
  { path: "/inscription", component: Register },
  { path: "/connexion", component: Login },
  { path: "/annonces", component: Annonces },
  { path: "/forums", component: Forums },
  { path: "/formations", component: Formations },
  { path: "/evenements", component: Evenements },
  {
    path: '/formations/:id',
    name: 'formation-detail',
    component: () => import('./pages/FormationDetail.vue')
  },
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
        path: "forums",
        name: "mes-forums",
        component: () => import("./pages/profilpages/Forums.vue"),
      },
      {
        path: "password",
        name: "modification-password",
        component: () => import("./pages/profilpages/ChangePassword.vue"),
      },
      {
        path: "modifyAnnonce/:id",
        name: "modification-annonce",
        component: () => import("./pages/profilpages/modifyAnnonce.vue"),
      },
      {
        path: "seeAnnonce/:id",
        name: "see-annonce",
        component: () => import("./pages/profilpages/SeeAnnonce.vue"),
      },
      {
        path: "createAnnonce",
        name: "create-annonce",
        component: () => import("./pages/profilpages/CreateAnnonce.vue"),
      },
      {
        path: "reserveCasier/:id",
        name: "reserve-casier",
        component: () => import("./pages/profilpages/ReserverCasier.vue"),
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
