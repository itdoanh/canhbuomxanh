(function honorPage() {
  if (window.gsap) {
    window.gsap
      .timeline({ defaults: { ease: 'power2.out' } })
      .from('.honor-head .kicker', { y: 18, opacity: 0, duration: 0.4 })
      .from('.honor-head h1', { y: 24, opacity: 0, duration: 0.55 }, '-=0.2')
      .from('.stats-strip .glass-card', { y: 20, opacity: 0, duration: 0.45, stagger: 0.09 }, '-=0.25')
      .from('.badge-grid .glass-card', { y: 20, opacity: 0, duration: 0.45, stagger: 0.1 }, '-=0.2');
  }

  const stats = {
    claimedCount: 248,
    votesCount: 19340,
    flagsCount: 112,
  };

  function animateStats() {
    Object.entries(stats).forEach(function (entry, index) {
      const key = entry[0];
      const value = entry[1];
      const el = document.getElementById(key);
      if (!el) return;

      setTimeout(function () {
        el.innerHTML = value;
      }, 320 + index * 180);
    });
  }

  function warmupApi() {
    if (!window.axios) return;
    window.axios
      .get('http://localhost:8080/api/v1/health', { timeout: 800 })
      .then(function () {
        return null;
      })
      .catch(function () {
        return null;
      });
  }

  animateStats();
  warmupApi();
})();
