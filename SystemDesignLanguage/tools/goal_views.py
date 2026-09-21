"""Goal traceability and feature slices, using explicit validated model facts only."""


def build_goal_views(views):
    relations, kinds = views.relations, views.kinds

    def names(kind):
        return sorted(n for n, k in kinds.items() if k == kind)

    def gap(viewpoint, name, code, message, mode=None, feature=None):
        views.gaps.append(dict(viewpoint=viewpoint, model_id=name, code=code,
                              message=message, mode=mode, feature=feature))

    if not names('usecase'):
        gap('VP01', '—', 'NO_USECASES', 'Ingen UseCase er deklarert.')
    for goal in names('usecase'):
        actors = [r for r in relations['pursues'] if r['object'] == goal]
        features = [r for r in relations['supports'] if r['object'] == goal]
        targets = {goal} | {r['subject'] for r in features}
        contributions = [r for r in relations['contributes-to'] if r['object'] in targets]
        direct = [r for r in contributions if r['object'] == goal]
        # Separate the one-hop goal map from feature contributions. This avoids
        # long bypass edges through several ranks in the supported renderer.
        views.diagram('VP01-' + goal, 'Bruksmål: ' + goal,
                      actors + features + direct, [goal])
        if not actors:
            gap('VP01', goal, 'NO_ACTOR', 'Ingen Actor pursues dette bruksmålet.')
        if not contributions:
            gap('VP01', goal, 'NO_CONTRIBUTION', 'Ingen direkte eller Feature-formidlede Functionality-bidrag.')
    for actor in names('actor'):
        if not any(r['subject'] == actor for r in relations['pursues']):
            gap('VP01', actor, 'NO_GOAL', 'Actor har ingen pursues-relasjon.')

    if not names('feature'):
        gap('VP07', '—', 'NO_FEATURES', 'Ingen Feature er deklarert.')
    for feature in names('feature'):
        contributions = [r for r in relations['contributes-to'] if r['object'] == feature]
        views.diagram('VP01-feature-' + feature, 'Functionality-bidrag til Feature: ' + feature,
                      contributions, [feature])
        functions = {r['subject'] for r in contributions}
        owners = [r for r in relations['owns'] if r['object'] in functions]
        allocations = [r for r in relations['allocated-to'] if r['subject'] in functions]
        modes = sorted({r['mode'] for r in allocations})
        if not any(r['subject'] == feature for r in relations['supports']):
            gap('VP01', feature, 'NO_GOAL', 'Feature støtter ingen deklarert UseCase.')
        if not contributions:
            gap('VP07', feature, 'NO_CONTRIBUTION', 'Feature har ingen Functionality-bidrag.')
        # An absent mode is never replaced with a guessed global/default context.
        for mode in modes or [None]:
            placed = [r for r in allocations if r['mode'] == mode]
            suffix = mode or 'unallocated'
            title = 'Feature: ' + feature + ' — ' + ('modus ' + mode if mode else 'uten allokering')
            views.diagram('VP07-' + feature + '-' + suffix, title,
                          contributions + owners + placed, [feature])
            located = {r['subject'] for r in placed}
            for function in sorted(functions - located):
                gap('VP07', function, 'UNSPECIFIED_ALLOCATION',
                    'Container-allokering er uspesifisert for bidrag til ' + feature + '.',
                    mode=mode, feature=feature)
