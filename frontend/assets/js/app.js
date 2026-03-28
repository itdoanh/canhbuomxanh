(function bootstrap() {
  if (window.Lenis) {
    const lenis = new window.Lenis({
      lerp: 0.09,
      wheelMultiplier: 0.9,
      smoothWheel: true,
    });

    function raf(time) {
      lenis.raf(time);
      requestAnimationFrame(raf);
    }
    requestAnimationFrame(raf);
  }

  if (window.VanillaTilt) {
    window.VanillaTilt.init(document.querySelectorAll('.tilt-card'), {
      max: 8,
      speed: 550,
      glare: true,
      'max-glare': 0.2,
      scale: 1.01,
    });
  }
})();
