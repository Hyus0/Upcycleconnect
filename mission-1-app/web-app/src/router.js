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
    ],
  },
];

export default createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (to.hash) {
      return {
        el: to.hash,
        behavior: "smooth",
      };
    }
    return { top: 0 };
  },
});
