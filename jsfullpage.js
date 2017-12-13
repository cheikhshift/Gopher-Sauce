  $(document).ready(function() {
      $('#fullpage').fullpage({scrollOverflow:true,
      anchors: ['start','overview','expsol','testtools','effic','usecases'] })
  });
       var gaJsHost = (("https:" == document.location.protocol) ? "https://ssl." : "http://www.");
            document.write(unescape("%3Cscript src='" + gaJsHost + "google-analytics.com/ga.js' type='text/javascript'%3E%3C/script%3E"));
            try {
              var pageTracker = _gat._getTracker("UA-71342160-1");
            pageTracker._trackPageview();
            } catch(err) {}