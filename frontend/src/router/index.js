import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue'),
    meta: { requiresAuth: true, scrapbook: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { guest: true, scrapbook: true }
  },
  {
    path: '/invite/:code',
    name: 'InviteRegister',
    component: () => import('../views/InviteRegister.vue'),
    meta: { guest: true, scrapbook: true }
  },
  {
    path: '/places',
    name: 'Places',
    component: () => import('../views/Places.vue'),
    meta: { requiresAuth: true, scrapbook: true }
  },
  {
    path: '/places/new',
    name: 'PlaceCreate',
    component: () => import('../views/PlaceForm.vue'),
    meta: { requiresAuth: true, scrapbook: true }
  },
  {
    path: '/places/:id',
    name: 'PlaceDetail',
    component: () => import('../views/PlaceDetail.vue'),
    meta: { requiresAuth: true, scrapbook: true }
  },
  {
    path: '/places/:id/edit',
    name: 'PlaceEdit',
    component: () => import('../views/PlaceForm.vue'),
    meta: { requiresAuth: true, scrapbook: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('../views/Profile.vue'),
    meta: { requiresAuth: true, scrapbook: true }
  },
  {
    path: '/map',
    name: 'Map',
    component: () => import('../views/MapPage.vue'),
    meta: { requiresAuth: true, scrapbook: true }
  },
  {
    path: '/invites',
    name: 'Invites',
    component: () => import('../views/Invites.vue'),
    meta: { requiresAuth: true, scrapbook: true }
  },
  {
    path: '/cities/:name',
    name: 'CityPage',
    component: () => import('../views/CityPage.vue'),
    meta: { requiresAuth: true, scrapbook: true }
  },
  {
    path: '/people/:id',
    name: 'PersonPage',
    component: () => import('../views/PersonPage.vue'),
    meta: { requiresAuth: true, scrapbook: true }
  },
  {
    path: '/gems',
    name: 'GemsHub',
    component: () => import('../views/GemsHub.vue'),
    meta: { requiresAuth: true, scrapbook: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')

  if (to.meta.requiresAuth && !token) {
    next({ name: 'Login' })
  } else if (to.meta.guest && token) {
    next({ name: 'Home' })
  } else {
    next()
  }
})

export default router
