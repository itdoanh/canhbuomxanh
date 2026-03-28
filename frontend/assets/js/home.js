(function homePageMotion() {
  if (!window.gsap) return;

  const tl = window.gsap.timeline({ defaults: { ease: 'power2.out' } });
  tl.from('.hero .kicker', { y: 24, opacity: 0, duration: 0.45 })
    .from('.hero h1', { y: 36, opacity: 0, duration: 0.6 }, '-=0.2')
    .from('.hero-copy', { y: 18, opacity: 0, duration: 0.5 }, '-=0.25')
    .from('.cta-row .btn', { y: 16, opacity: 0, duration: 0.45, stagger: 0.08 }, '-=0.2')
    .from('.hero-grid .glass-card', { y: 22, opacity: 0, duration: 0.5, stagger: 0.1 }, '-=0.2');
})();
