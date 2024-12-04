import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Hammer",
  description: "Screw deployments",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Getting Started', link: '/getting-started' }
    ],

    sidebar: [
      {
        text: "Guides",
        items: [
          { text: "Introduction", link: "/introduction" },
          { text: "Getting Started", link: "/getting-started" }
        ]
      },
      {
        text: 'Reference',
        items: [
          { text: 'Configuration Options', link: '/reference/configuration' },
          { text: 'CLI Options', link: '/reference/cli-options' }
        ]
      }
    ],
    search: {
      provider: "local"
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/matfire/hammer' }
    ]
  }
})
